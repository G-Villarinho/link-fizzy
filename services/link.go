package services

import (
	"context"
	"fmt"
	"time"

	"github.com/g-villarinho/link-fizz-api/config"
	"github.com/g-villarinho/link-fizz-api/models"
	"github.com/g-villarinho/link-fizz-api/pkgs/di"
	"github.com/g-villarinho/link-fizz-api/repositories"
	"github.com/google/uuid"
)

const shortCodeLength = 8

type LinkService interface {
	CreateLink(ctx context.Context, originalURL string) error
	GetOriginalURLByShortCode(ctx context.Context, shortCode string) (string, error)
	GetShortURLs(ctx context.Context) ([]string, error)
}

type linkService struct {
	i  *di.Injector
	us UtilsService
	lr repositories.LinkRepository
}

func NewLinkService(i *di.Injector) (LinkService, error) {
	utilsService, err := di.Invoke[UtilsService](i)
	if err != nil {
		return nil, fmt.Errorf("invoke services.UtilsService: %w", err)
	}

	linkRepository, err := di.Invoke[repositories.LinkRepository](i)
	if err != nil {
		return nil, fmt.Errorf("invoke repositories.LinkRepository: %w", err)
	}

	return &linkService{
		i:  i,
		us: utilsService,
		lr: linkRepository,
	}, nil
}

func (l *linkService) CreateLink(ctx context.Context, originalURL string) error {
	shortCode, err := l.us.GenerateShortCode(shortCodeLength)
	if err != nil {
		return fmt.Errorf("generate short code: %w", err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("generate UUID: %w", err)
	}

	link := models.Link{
		ID:          id.String(),
		OriginalURL: originalURL,
		ShortCode:   shortCode,
		CreatedAt:   time.Now().UTC(),
	}

	if err := l.lr.CreateLink(ctx, link); err != nil {
		return fmt.Errorf("create link: %w", err)
	}

	return nil
}

func (l *linkService) GetOriginalURLByShortCode(ctx context.Context, shortCode string) (string, error) {
	originalUrl, err := l.lr.GetOriginalURLByShortCode(ctx, shortCode)
	if err != nil {
		return "", fmt.Errorf("get original URL by short code: %w", err)
	}

	if originalUrl == "" {
		return "", models.ErrLinkNotFound
	}

	return originalUrl, nil
}

func (l *linkService) GetShortURLs(ctx context.Context) ([]string, error) {
	shortCodes, err := l.lr.GetAllShortCodes(ctx)
	if err != nil {
		return nil, fmt.Errorf("get all short codes: %w", err)
	}

	if len(shortCodes) == 0 {
		return []string{}, nil
	}

	shortUrls := make([]string, len(shortCodes))
	apiURL := config.Env.APIURL

	for i, shortCode := range shortCodes {
		shortUrls[i] = fmt.Sprintf("%s/%s", apiURL, shortCode)
	}

	return shortUrls, nil
}
