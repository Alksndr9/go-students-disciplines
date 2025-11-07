package logger

import "go.uber.org/zap"

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func InitLogger(env string) *zap.Logger {
	var log *zap.Logger

	switch env {
	case envLocal:
		log = zap.Must(zap.NewDevelopment())
	case envDev:
		log = zap.Must(zap.NewDevelopment())
	case envProd:
		log = zap.Must(zap.NewProduction())
	}

	return log
}
