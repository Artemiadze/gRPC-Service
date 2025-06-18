package app

import (
	"time"

	grpcapp "github.com/Artemiadze/gRPC-Service/internal/app/grpc"
	postgres "github.com/Artemiadze/gRPC-Service/internal/repository"
	"github.com/Artemiadze/gRPC-Service/internal/services"
	"go.uber.org/zap"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *zap.Logger,
	grpcPort int,
	dsn string,
	tokenTTL time.Duration,
) *App {
	// Инициализация хранилища
	storage, err := postgres.New(dsn)
	if err != nil {
		panic(err)
	}

	authService := services.New(log, storage, storage, storage, tokenTTL)

	// инициализация gRPC сервера
	grpcApp := grpcapp.New(log, authService, grpcPort)
	return &App{
		GRPCServer: grpcApp,
	}
}
