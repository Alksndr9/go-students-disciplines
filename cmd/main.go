package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"gitgub.com/Alksndr9/go-students-disciplines/internal/config"
	"gitgub.com/Alksndr9/go-students-disciplines/internal/modules"
	"gitgub.com/Alksndr9/go-students-disciplines/internal/responder"
	"gitgub.com/Alksndr9/go-students-disciplines/internal/router"
	"gitgub.com/Alksndr9/go-students-disciplines/pkg/logger"
	"gitgub.com/Alksndr9/go-students-disciplines/pkg/psql"
	"go.uber.org/zap"
)

const (
	_shutdownPeriod      = 15 * time.Second
	_shutdownHardPeriod  = 3 * time.Second
	_readinessDrainDelay = 5 * time.Second
)

var isShuttingDown atomic.Bool

func main() {
	cfg := config.MustLoad()

	logger := logger.InitLogger(cfg.Env)

	rootCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	db, err := psql.Connect(rootCtx, cfg, logger)
	if err != nil {
		logger.Error("failed to init db", zap.Error(err))
		os.Exit(1)
	}
	logger.Info("successfully connected to the db")

	ongoingCtx, stopOngoingGracefully := context.WithCancel(context.Background())

	logger.Info("starting students-disciplines", zap.String("env", cfg.Env))

	responder := responder.NewResponder(logger)

	storages := modules.NewStorages(db)
	controllers := modules.NewControllers(responder, logger, storages)

	router := router.NewRouter(controllers)

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
		BaseContext: func(_ net.Listener) context.Context {
			return ongoingCtx
		},
	}

	go func() {
		logger.Info("starting server", zap.String("address", cfg.Address))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	<-rootCtx.Done()
	stop()
	isShuttingDown.Store(true)
	logger.Info("Received shutdown signal, shutting down.")

	time.Sleep(_readinessDrainDelay)
	logger.Info("Readiness check propagated, now waiting for ongoing requests to finish.")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), _shutdownPeriod)
	defer cancel()

	err = srv.Shutdown(shutdownCtx)
	stopOngoingGracefully()

	if err != nil {
		logger.Warn("Failed to wait for ongoing requests to finish, waiting for forced cancellation.")
		time.Sleep(_shutdownHardPeriod)
	}

	logger.Info("Server shut down gracefully.")

	// TO-DO: Users CRUD
	// TO-DO: Users usecases -> validation
	// TO-DO: Users service
}
