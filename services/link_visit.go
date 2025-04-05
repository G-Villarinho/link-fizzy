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

type LinkVisitService interface {
	CreateLinkVisit(ctx context.Context, linkID, ipAddress, userAgent string) error
}

type linkVisitService struct {
	i   *di.Injector
	lvr repositories.LinkVisitRepository
}

func NewLinkVisitService(i *di.Injector) (LinkVisitService, error) {
	linkVisitRepository, err := di.Invoke[repositories.LinkVisitRepository](i)
	if err != nil {
		return nil, fmt.Errorf("invoke repositories.LinkVisitRepository: %w", err)
	}

	return &linkVisitService{
		i:   i,
		lvr: linkVisitRepository,
	}, nil
}

func (l *linkVisitService) CreateLinkVisit(ctx context.Context, linkID string, ipAddress string, userAgent string) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("generate uuid: %w", err)
	}

	linkVisit := &models.LinkVisit{
		ID:        id.String(),
		LinkID:    linkID,
		IP:        ipAddress,
		Agent:     userAgent,
		VisitedAt: time.Now().UTC(),
	}

	if err := l.lvr.CreateLinkVisit(ctx, linkVisit); err != nil {
		return fmt.Errorf("create link visit: %w", err)
	}

	return nil
}
