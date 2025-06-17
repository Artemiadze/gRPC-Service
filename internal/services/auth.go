package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	err_internal "github.com/Artemiadze/gRPC-Service/internal/errors"
	"github.com/Artemiadze/gRPC-Service/internal/lib/jwt"
	"github.com/Artemiadze/gRPC-Service/internal/models"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
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

func (a *AuthService) Login(ctx context.Context, email string, password string, appID int64) (string, error) {
	const op = "AuthService.Login"
	log := a.log.With(zap.String("method", op), zap.String("email", email))

	log.Info("attempting to login user")
	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, err_internal.ErrUserNotFound) {
			a.log.Error("user not found", zap.Error(err))
			return "", fmt.Errorf("user not found: %w", err)
		}
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Error("password mismatch", zap.Error(err))
		return "", fmt.Errorf("password mismatch: %w", err)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("failed to get app: %s %w", op, err)
	}

	log.Info("user logged in successfully")

	token, err := jwt.GenerateToken(user, app, a.tokenTTL)
	if err != nil {
		return "", fmt.Errorf("failed to get app: %s %w", op, err)
	}

	return token, nil
}

func (a *AuthService) RegisterNewUser(ctx context.Context, email string, password string, appID int) (int64, error) {
	const op = "AuthService.RegisterNewUser"
	log := a.log.With(zap.String("method", op), zap.String("email", email))

	log.Info("registering new user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to hash password", zap.Error(err))
		return 0, err
	}

	id, err := a.usrSaver.SaveUser(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save user", zap.Error(err))
		return 0, err
	}

	log.Info("user registered successfully", zap.Int64("userID", id))
	return id, nil
}

func (a *AuthService) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "AuthService.IsAdmin"
	log := a.log.With(zap.String("method", op), zap.Int64("ID", userID))

	log.Info("checking if user is admin")

	isAdmin, err := a.usrProvider.IsAdmin(ctx, userID)
	if err != nil {
		return false, fmt.Errorf("failed to check admin status: %s %w", op, err)
	}

	log.Info("checked admin status successfully", zap.Bool("isAdmin", isAdmin))
	return isAdmin, nil
}
