package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/g-villarinho/link-fizz-api/models"
	"github.com/g-villarinho/link-fizz-api/pkgs/di"
)

type SessionRepository interface {
	CreateSession(ctx context.Context, session *models.Session) error
	DeleteSession(ctx context.Context, sessionID string) error
}

type sessionRepository struct {
	i  *di.Injector
	db *sql.DB
}

func NewSessionRepository(i *di.Injector) (SessionRepository, error) {
	db, err := di.Invoke[*sql.DB](i)
	if err != nil {
		return nil, fmt.Errorf("invoke db: %w", err)
	}

	return &sessionRepository{
		i:  i,
		db: db,
	}, nil
}

func (s *sessionRepository) CreateSession(ctx context.Context, session *models.Session) error {
	statment, err := s.db.PrepareContext(ctx, "INSERT INTO sessions (id, user_id, token, ip, agent, create_at, expire_at) VALUES (?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("prepare statement: %w", err)
	}
	defer statment.Close()

	_, err = statment.ExecContext(ctx, session.ID, session.UserID, session.Token, session.IP, session.Agent, session.CreateAt, session.ExpireAt)
	if err != nil {
		return fmt.Errorf("exec statement: %w", err)
	}

	return nil
}

func (s *sessionRepository) DeleteSession(ctx context.Context, sessionID string) error {
	statment, err := s.db.PrepareContext(ctx, "DELETE FROM sessions WHERE id = ?")
	if err != nil {
		return fmt.Errorf("prepare statement: %w", err)
	}
	defer statment.Close()

	_, err = statment.ExecContext(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("exec statement: %w", err)
	}

	return nil
}
