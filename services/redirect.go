package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/g-villarinho/link-fizz-api/pkgs/di"
)

type RedirectService interface {
	GetOriginalURLWithTracking(ctx context.Context, linkShortCode, userAgent, ipAddress string) (string, error)
}

type redirectService struct {
	i   *di.Injector
	ls  LinkService
	lvs LinkVisitService
}

func NewRedirectService(i *di.Injector) (RedirectService, error) {
	linkService, err := di.Invoke[LinkService](i)
	if err != nil {
		return nil, fmt.Errorf("invoke services.Link: %w", err)
	}

	linkVisitService, err := di.Invoke[LinkVisitService](i)
	if err != nil {
		return nil, fmt.Errorf("invoke services.LinkVisit: %w", err)
	}

	return &redirectService{
		i:   i,
		ls:  linkService,
		lvs: linkVisitService,
	}, nil
}

func (r *redirectService) GetOriginalURLWithTracking(ctx context.Context, linkShortCode string, userAgent string, ipAddress string) (string, error) {
	logger := slog.With(
		"service", "redirect",
		"method", "GetOriginalURLWithTracking",
	)

	link, err := r.ls.GetLinkByShortCode(ctx, linkShortCode)
	if err != nil {
		return "", fmt.Errorf("get original URL by short code: %w", err)
	}

	go func() {
		if err = r.lvs.CreateLinkVisit(context.Background(), link.ID, ipAddress, userAgent); err != nil {
			logger.Error("create link visit", "error", err)
			return
		}
	}()

	return link.OriginalURL, nil
}
