package main

import (
	"github.com/Artemiadze/gRPC-Service/internal/app"
	"github.com/Artemiadze/gRPC-Service/internal/config"
	"go.uber.org/zap"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	// для запуска go run cmd/sso/main.go --config=./config/local.yaml

	// загрузка конфигурации
	cfg := config.MustLoad()

	// инициализация логгера
	logger := setupLogger(cfg.Env)
	defer logger.Sync() // flush буфера

	logger.Info("starting service",
		zap.String("env", cfg.Env),
		zap.String("storage_path", cfg.StoragePath),
	)
	//logger.Debug("Debug message")

	// инициализация приложения (app)
	application := app.New(logger, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	application.GRPCServer.MustRun() // запускаем gRPC-сервер приложения
	// запустить gRPC-сервер приложения
}

func setupLogger(env string) *zap.Logger {
	var logger *zap.Logger
	var err error

	switch env {
	case envLocal:
		cfg := zap.NewDevelopmentConfig()
		logger, err = cfg.Build()
	case envDev:
		cfg := zap.NewDevelopmentConfig()
		cfg.Encoding = "json"
		logger, err = cfg.Build()
	case envProd:
		cfg := zap.NewProductionConfig()
		logger, err = cfg.Build()
	default:
		logger = zap.NewExample() // fallback logger
	}

	if err != nil {
		// В случае ошибки инициализации логгера — паника
		panic("cannot initialize zap logger: " + err.Error())
	}

	return logger
}
