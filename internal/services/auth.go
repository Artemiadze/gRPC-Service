package services

import (
	"context"
	"time"

	"github.com/Artemiadze/gRPC-Service/internal/models"
	"go.uber.org/zap"
)

type AuthService struct {
	// Add any dependencies or configurations needed for the AuthService
	log         *zap.Logger
	usrSaver    Storage
	usrProvider Storage
	appProvider Storage
	tokenTTL    time.Duration
}

type Storage interface {
	// Define methods that the storage layer should implement
	User(ctx context.Context, email string) (user models.User, err error)
	IsAdmin(ctx context.Context, uid int64) (isAdmin bool, err error)
	SaveUser(ctx context.Context, email string, passHash []byte) (uid int64, err error)
	App(ctx context.Context, uid int64) (app models.App, err error)
}

// New creates a new instance of AuthService with the provided dependencies.
func New(
	log *zap.Logger,
	userSaver Storage,
	userProvider Storage,
	appProvider Storage,
	tokenTTL time.Duration,
) *AuthService {
	return &AuthService{
		usrSaver:    userSaver,
		usrProvider: userProvider,
		log:         log,
		appProvider: appProvider,
		tokenTTL:    tokenTTL,
	}
}

func (a *AuthService) Login(ctx context.Context, email string, password string, appID int) (string, error) {
	panic("not implemented yet") // TODO: Implement GetUser method
}

func (a *AuthService) RegisterNewUser(ctx context.Context, email string, password string, appID int) (string, error) {
	panic("not implemented yet") // TODO: Implement GetUser method
}

func (a *AuthService) IsAdmin(ctx context.Context, email string, password string, appID int) (string, error) {
	panic("not implemented yet") // TODO: Implement GetUser method
}
