package auth

import (
	"context"

	ssov1 "github.com/Artemiadze/gRPC-Service/gen/go/sso"
	"google.golang.org/grpc"
)

type serverAPI struct {
	ssov1.UnimplementedAuthServer
	auth Auth
}

// Интерфейс, который мы передавали в grpcApp
type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int,
	) (token string, err error)
	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID int64, err error)
}

func Register(gRPCServer *grpc.Server) {
	ssov1.RegisterAuthServer(gRPCServer, &serverAPI{})
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *ssov1.LoginRequest,
) (*ssov1.LoginResponse, error) {
	//panic("implemented me")
	return &ssov1.LoginResponse{
		Token: req.GetEmail(),
	}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *ssov1.RegisterRequest,
) (*ssov1.RegisterResponse, error) {
	panic("implemented me")
}

func (s *serverAPI) IsAdmin(
	ctx context.Context,
	req *ssov1.IsAdminRequest,
) (*ssov1.IsAdminResponse, error) {
	panic("implemented me")
}
