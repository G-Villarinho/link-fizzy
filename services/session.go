package services

import (
	"context"
	"fmt"
	"time"

	"github.com/g-villarinho/link-fizz-api/models"
	"github.com/g-villarinho/link-fizz-api/pkgs/di"
	"github.com/g-villarinho/link-fizz-api/repositories"
	"github.com/google/uuid"
)

type SessionService interface {
	CreateSession(ctx context.Context, userID, IPAddress, userAgent string) (string, error)
	DeleteSession(ctx context.Context, sessionID string) error
}

type sessionService struct {
	i  *di.Injector
	ts TokenService
	sr repositories.SessionRepository
}

func NewSessionService(i *di.Injector) (SessionService, error) {
	tokenService, err := di.Invoke[TokenService](i)
	if err != nil {
		return nil, fmt.Errorf("invoke services.Token: %w", err)
	}

	sessionRepository, err := di.Invoke[repositories.SessionRepository](i)
	if err != nil {
		return nil, fmt.Errorf("invoke repositories.Session: %w", err)
	}

	return &sessionService{
		i:  i,
		ts: tokenService,
		sr: sessionRepository,
	}, nil
}

func (s *sessionService) CreateSession(ctx context.Context, userID string, IPAddress string, userAgent string) (string, error) {
	sessionID, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("generate session id: %w", err)
	}

	now := time.Now().UTC()
	expiresAt := now.Add(time.Hour * 24 * 7)

	token, err := s.ts.GenerateToken(ctx, userID, sessionID.String(), now, expiresAt)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	session := &models.Session{
		ID:       sessionID.String(),
		UserID:   userID,
		Token:    token,
		IP:       IPAddress,
		Agent:    userAgent,
		CreateAt: now,
		ExpireAt: expiresAt,
	}

	if err := s.sr.CreateSession(ctx, session); err != nil {
		return "", fmt.Errorf("create session: %w", err)
	}

	return token, nil
}

func (s *sessionService) DeleteSession(ctx context.Context, sessionID string) error {
	if err := s.sr.DeleteSession(ctx, sessionID); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}
