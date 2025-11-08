package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"gitgub.com/Alksndr9/go-students-disciplines/internal/config"
	"gitgub.com/Alksndr9/go-students-disciplines/internal/router"
	"gitgub.com/Alksndr9/go-students-disciplines/pkg/logger"
	"gitgub.com/Alksndr9/go-students-disciplines/pkg/psql"
	"go.uber.org/zap"
)

func main() {
	cfg := config.MustLoad()

	logger := logger.InitLogger(cfg.Env)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	_, err := psql.Connect(ctx, cfg, logger)
	if err != nil {
		logger.Error("failed to init db", zap.Error(err))
		os.Exit(1)
	}
	logger.Info("successfully connected to the db")

	logger.Info("starting students-disciplines", zap.String("env", cfg.Env))

	router := router.NewRouter()
	logger.Info("starting server", zap.String("address", cfg.Address))

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.Timeout,
		WriteTimeout: cfg.Timeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	if err = srv.ListenAndServe(); err != nil {
		logger.Error("failde to start server")
	}

	// TO-DO: Users CRUD

	// TO-DO: Gracefull-Shutdown

	// TO-DO: router
}
