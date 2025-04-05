package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/g-villarinho/link-fizz-api/models"
	"github.com/g-villarinho/link-fizz-api/pkgs/di"
)

type LogoutRepository interface {
	CreateLogout(ctx context.Context, logout *models.Logout) error
	GetLogoutByToken(ctx context.Context, token string) (*models.Logout, error)
}

type logoutRepository struct {
	i  *di.Injector
	db *sql.DB
}

func NewLogoutRepository(i *di.Injector) (LogoutRepository, error) {
	db, err := di.Invoke[*sql.DB](i)
	if err != nil {
		return nil, fmt.Errorf("invoke db: %w", err)
	}

	return &logoutRepository{
		i:  i,
		db: db,
	}, nil
}

func (l *logoutRepository) CreateLogout(ctx context.Context, logout *models.Logout) error {
	statement, err := l.db.PrepareContext(ctx, "INSERT INTO logouts (id, token, user_id, revoked_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("prepare insert: %w", err)
	}
	defer statement.Close()

	_, err = statement.ExecContext(ctx, logout.ID, logout.Token, logout.UserID, logout.RevokedAt)
	if err != nil {
		return fmt.Errorf("exec insert: %w", err)
	}

	return nil
}

func (l *logoutRepository) GetLogoutByToken(ctx context.Context, token string) (*models.Logout, error) {
	statement, err := l.db.PrepareContext(ctx, "SELECT id, token, user_id, revoked_at FROM logouts WHERE token = ?")
	if err != nil {
		return nil, fmt.Errorf("prepare select: %w", err)
	}
	defer statement.Close()

	logout := &models.Logout{}
	err = statement.QueryRowContext(ctx, token).Scan(&logout.ID, &logout.Token, &logout.UserID, &logout.RevokedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("query row: %w", err)
	}

	return logout, nil
}
