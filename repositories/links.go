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
	statement, err := l.db.
		PrepareContext(ctx, "INSERT INTO links (id, original_url, short_code, created_at) VALUES (?, ?, ?, ?)")

	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.ExecContext(ctx, link.ID, link.OriginalURL, link.ShortCode, link.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (l *linkRepository) GetOriginalURLByShortCode(ctx context.Context, shortCode string) (string, error) {
	statement, err := l.db.PrepareContext(ctx, "SELECT original_url FROM links WHERE short_code = ?")
	if err != nil {
		return "", err
	}
	defer statement.Close()

	var originalURL string
	err = statement.QueryRowContext(ctx, shortCode).Scan(&originalURL)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}

		return "", err
	}

	return originalURL, nil
}

func (l *linkRepository) GetLinkByID(ctx context.Context, ID string) (*models.Link, error) {
	statement, err := l.db.
		PrepareContext(ctx, "SELECT id, original_url, short_code, created_at, updated_at FROM links WHERE id = ?")

	if err != nil {
		return nil, err
	}
	defer statement.Close()

	var link models.Link
	err = statement.
		QueryRowContext(ctx, ID).
		Scan(&link.ID, &link.OriginalURL, &link.ShortCode, &link.CreatedAt, &link.UpdatedAt)

	if err != nil {
		return nil, err
	}

	return &link, nil
}
