package services

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/g-villarinho/link-fizz-api/models"
	"github.com/g-villarinho/link-fizz-api/pkgs/di"
	"github.com/g-villarinho/link-fizz-api/repositories"
	"github.com/google/uuid"
)

type LogoutService interface {
	CreateLogout(ctx context.Context, userID *string, token string) error
	IsLogoutRegistered(ctx context.Context, token string) (bool, error)
}

type logoutService struct {
	i  *di.Injector
	lr repositories.LogoutRepository
}

func NewLogoutService(i *di.Injector) (LogoutService, error) {
	logoutRepository, err := di.Invoke[repositories.LogoutRepository](i)
	if err != nil {
		return nil, fmt.Errorf("invoke repositories.LogoutRepository: %w", err)
	}

	return &logoutService{
		i:  i,
		lr: logoutRepository,
	}, nil
}

func (l *logoutService) CreateLogout(ctx context.Context, userID *string, token string) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("uuid.NewRandom: %w", err)
	}

	logout := &models.Logout{
		ID:        id.String(),
		Token:     token,
		RevokedAt: time.Now().UTC(),
	}

	if userID != nil {
		logout.UserID = sql.NullString{String: *userID, Valid: true}
	}

	if err := l.lr.CreateLogout(ctx, logout); err != nil {
		return fmt.Errorf("create logout: %w", err)
	}

	return nil
}

func (l *logoutService) IsLogoutRegistered(ctx context.Context, token string) (bool, error) {
	logoutFromToken, err := l.lr.GetLogoutByToken(ctx, token)
	if err != nil {
		return false, fmt.Errorf("get logout by token: %w", err)
	}

	return logoutFromToken != nil, nil
}
