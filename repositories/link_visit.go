package repositories

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/g-villarinho/link-fizz-api/models"
	"github.com/g-villarinho/link-fizz-api/pkgs/di"
)

type LinkVisitRepository interface {
	CreateLinkVisit(ctx context.Context, linkVisit *models.LinkVisit) error
}

type linkVisitRepository struct {
	i  *di.Injector
	db *sql.DB
}

func NewLinkVisitRepository(i *di.Injector) (LinkVisitRepository, error) {
	db, err := di.Invoke[*sql.DB](i)
	if err != nil {
		return nil, fmt.Errorf("invoke sql.db: %w", err)
	}

	return &linkVisitRepository{
		i:  i,
		db: db,
	}, nil
}

func (l *linkVisitRepository) CreateLinkVisit(ctx context.Context, linkVisit *models.LinkVisit) error {
	statement, err := l.db.PrepareContext(ctx, "INSERT INTO link_visits (id, link_id, ip, agent, visited_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return fmt.Errorf("prepare select: %w", err)
	}
	defer statement.Close()

	_, err = statement.ExecContext(ctx, linkVisit.ID, linkVisit.LinkID, linkVisit.IP, linkVisit.Agent, linkVisit.VisitedAt)
	if err != nil {
		return fmt.Errorf("exec select: %w", err)
	}

	return nil
}
