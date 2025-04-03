package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/g-villarinho/link-fizz-api/models"
	"github.com/g-villarinho/link-fizz-api/pkgs/di"
)

type LinkRepository interface {
	CreateLink(ctx context.Context, link models.Link) error
	GetOriginalURLByShortCode(ctx context.Context, shortCode string) (string, error)
	GetLinkByID(ctx context.Context, ID string) (*models.Link, error)
	GetAllShortCodesByUserID(ctx context.Context, userID string) ([]string, error)
}

type linkRepository struct {
	i  *di.Injector
	db *sql.DB
}

func NewLinkRepository(i *di.Injector) (LinkRepository, error) {
	db, err := di.Invoke[*sql.DB](i)
	if err != nil {
		return nil, fmt.Errorf("invoke sql.DB: %w", err)
	}

	return &linkRepository{
		i:  i,
		db: db,
	}, nil
}

func (l *linkRepository) CreateLink(ctx context.Context, link models.Link) error {
	statement, err := l.db.PrepareContext(ctx, "INSERT INTO links (id, original_url, short_code, created_at) VALUES (?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("prepare insert: %w", err)
	}
	defer statement.Close()

	_, err = statement.ExecContext(ctx, link.ID, link.OriginalURL, link.ShortCode, link.CreatedAt)
	if err != nil {
		return fmt.Errorf("execute insert: %w", err)
	}

	return nil
}

func (l *linkRepository) GetOriginalURLByShortCode(ctx context.Context, shortCode string) (string, error) {
	var originalURL string
	err := l.db.QueryRowContext(ctx, "SELECT original_url FROM links WHERE short_code = ?", shortCode).Scan(&originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", fmt.Errorf("query original_url: %w", err)
	}
	return originalURL, nil
}

func (l *linkRepository) GetLinkByID(ctx context.Context, ID string) (*models.Link, error) {
	statement, err := l.db.PrepareContext(ctx, "SELECT id, original_url, short_code, created_at, updated_at FROM links WHERE id = ?")
	if err != nil {
		return nil, fmt.Errorf("prepare select: %w", err)
	}
	defer statement.Close()

	return l.scanLink(statement.QueryRowContext(ctx, ID))
}

func (l *linkRepository) scanLink(row *sql.Row) (*models.Link, error) {
	var link models.Link
	err := row.Scan(&link.ID, &link.OriginalURL, &link.ShortCode, &link.CreatedAt, &link.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("scan link: %w", err)
	}
	return &link, nil
}

func (l *linkRepository) GetAllShortCodesByUserID(ctx context.Context, userID string) ([]string, error) {
	statement, err := l.db.PrepareContext(ctx, "SELECT short_code FROM links WHERE user_id = ?")
	if err != nil {
		return nil, fmt.Errorf("prepare select: %w", err)
	}
	defer statement.Close()

	rows, err := statement.QueryContext(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("query select: %w", err)
	}
	defer rows.Close()

	var shortCodes []string
	for rows.Next() {
		var shortCode string
		if err := rows.Scan(&shortCode); err != nil {
			return nil, fmt.Errorf("scan short code: %w", err)
		}
		shortCodes = append(shortCodes, shortCode)
	}

	return shortCodes, nil
}
