package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"gitgub.com/Alksndr9/go-students-disciplines/internal/config"
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

	// TO-DO: router gin
}
