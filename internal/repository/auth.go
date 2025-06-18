package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_error "github.com/Artemiadze/gRPC-Service/internal/errors"
	"github.com/Artemiadze/gRPC-Service/internal/models"

	"github.com/lib/pq"
)

type repository struct {
	db *sql.DB
}

func New(dsn string) (*repository, error) {
	const op = "repository.postgres.New"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Проверим соединение
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("%s: ping error: %w", op, err)
	}

	return &repository{db: db}, nil
}

func (s *repository) Stop() error {
	return s.db.Close()
}

func (s *repository) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "repository.postgres.SaveUser"

	stmt, err := s.db.PrepareContext(ctx,
		`INSERT INTO users(email, pass_hash) VALUES($1, $2) RETURNING id`)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var id int64
	err = stmt.QueryRowContext(ctx, email, passHash).Scan(&id)
	if err != nil {
		var pgErr *pq.Error
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, _error.ErrUserExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *repository) User(ctx context.Context, email string) (models.User, error) {
	const op = "repository.postgres.User"

	stmt, err := s.db.PrepareContext(ctx,
		`SELECT id, email, pass_hash FROM users WHERE email = $1`)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var user models.User
	err = stmt.QueryRowContext(ctx, email).Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, _error.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *repository) App(ctx context.Context, id int) (models.App, error) {
	const op = "repository.postgres.App"

	stmt, err := s.db.PrepareContext(ctx,
		`SELECT id, name, secret FROM apps WHERE id = $1`)
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var app models.App
	err = stmt.QueryRowContext(ctx, id).Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, _error.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}

func (s *repository) IsAdmin(ctx context.Context, userID int64) (bool, error) {
	const op = "repository.postgres.IsAdmin"

	stmt, err := s.db.PrepareContext(ctx,
		`SELECT is_admin FROM users WHERE id = $1`)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var isAdmin bool
	err = stmt.QueryRowContext(ctx, userID).Scan(&isAdmin)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, fmt.Errorf("%s: %w", op, _error.ErrUserNotFound)
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return isAdmin, nil
}
