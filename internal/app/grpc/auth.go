package grpcapp

import (
	"fmt"
	"net"

	"go.uber.org/zap"

	authgrpc "github.com/Artemiadze/gRPC-Service/internal/grpc/Auth"

	"google.golang.org/grpc"
)

type App struct {
	log        *zap.Logger
	gRPCServer *grpc.Server
	port       int
}

// Create a new gRPC server application
func New(
	log *zap.Logger,
	port int,
) *App {
	gRPCServer := grpc.NewServer()

	authgrpc.Register(gRPCServer)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

// MustRun runs gRPC server and panics if any error occurs.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run runs gRPC server.
func (a *App) Run() error {
	const app = "grpcapp.Run"

	log := a.log.With(
		zap.String("app", app),
		zap.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", app, err)
	}

	log.Info("Starting gRPC server", zap.String("address", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", app, err)
	}

	return nil
}

// Stop stops gRPC server.
func (a *App) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(zap.String("op", op)).
		Info("stopping gRPC server", zap.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
