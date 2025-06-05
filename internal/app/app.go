package app

import (
	"time"

	"go.uber.org/zap"

	grpcapp "github.com/Artemiadze/gRPC-Service/internal/app/grpc"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *zap.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	// Инициализация хранилища

	// инициализация gRPC сервера
	grpcApp := grpcapp.New(log, grpcPort)
	return &App{
		GRPCServer: grpcApp,
	}
}
