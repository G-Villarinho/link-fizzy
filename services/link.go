package services

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/g-villarinho/link-fizz-api/config"
	"github.com/g-villarinho/link-fizz-api/models"
	"github.com/g-villarinho/link-fizz-api/pkgs/di"
	"github.com/g-villarinho/link-fizz-api/repositories"
	"github.com/google/uuid"
)

const shortCodeLength = 8

type LinkService interface {
	CreateLink(ctx context.Context, userID, destinationURL string, title, customCode *string) (*models.CreateLinkResponse, error)
	GetOriginalURLByShortCode(ctx context.Context, shortCode string) (string, error)
	GetUsersShortURLs(ctx context.Context, userID string) ([]string, error)
	GetLinkByShortCode(ctx context.Context, shortCode string) (*models.Link, error)
	GetLinksByUserID(ctx context.Context, userID string) ([]models.LinkResponse, error)
	GetLinkDetails(ctx context.Context, userID string, shortCode string) (*models.LinkResponse, error)
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

func (l *linkService) CreateLink(ctx context.Context, userID, destinationURL string, title, customCode *string) (*models.CreateLinkResponse, error) {
	var shortCode string

	if customCode != nil && *customCode != "" {
		cleanCode := strings.ReplaceAll(*customCode, " ", "")
		cleanCode = strings.ToLower(cleanCode)

		linkFromCode, err := l.lr.GetLinkByShortCode(ctx, cleanCode)
		if err != nil {
			return nil, fmt.Errorf("get link by short code: %w", err)
		}

		if linkFromCode != nil {
			return nil, models.ErrCustomCodeAlreadyExists
		}

		shortCode = cleanCode
	} else {
		generatedCode, err := l.us.GenerateShortCode(shortCodeLength)
		if err != nil {
			return nil, fmt.Errorf("generate short code: %w", err)
		}

		shortCode = generatedCode
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("generate UUID: %w", err)
	}

	link := models.Link{
		ID:          id.String(),
		Title:       sql.NullString{String: *title, Valid: title != nil},
		OriginalURL: destinationURL,
		ShortCode:   shortCode,
		UserID:      userID,
		CreatedAt:   time.Now().UTC(),
	}

	if err := l.lr.CreateLink(ctx, link); err != nil {
		return nil, fmt.Errorf("create link: %w", err)
	}

	return &models.CreateLinkResponse{
		ShortCode: shortCode,
	}, nil
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

func (l *linkService) GetUsersShortURLs(ctx context.Context, userID string) ([]string, error) {
	shortCodes, err := l.lr.GetAllShortCodesByUserID(ctx, userID)
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

func (l *linkService) GetLinkByShortCode(ctx context.Context, shortCode string) (*models.Link, error) {
	link, err := l.lr.GetLinkByShortCode(ctx, shortCode)
	if err != nil {
		return nil, fmt.Errorf("get link by short code: %w", err)
	}

	if link == nil {
		return nil, models.ErrLinkNotFound
	}

	return link, nil
}

func (l *linkService) GetLinksByUserID(ctx context.Context, userID string) ([]models.LinkResponse, error) {
	links, err := l.lr.GetLinksByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get links by user ID: %w", err)
	}

	if len(links) == 0 {
		return []models.LinkResponse{}, nil
	}

	apiURL := config.Env.APIURL
	linkResponses := make([]models.LinkResponse, len(links))

	for i, link := range links {
		linkResponses[i] = link.ToResponse(apiURL)
	}

	return linkResponses, nil
}

func (l *linkService) GetLinkDetails(ctx context.Context, userID string, shortCode string) (*models.LinkResponse, error) {
	link, err := l.lr.GetLinkByShortCode(ctx, shortCode)
	if err != nil {
		return nil, fmt.Errorf("get link by short code: %w", err)
	}

	if link == nil {
		return nil, models.ErrLinkNotFound
	}

	if link.UserID != userID {
		return nil, models.ErrLinkNotBelongToUser
	}

	response := link.ToResponse(config.Env.APIURL)
	return &response, nil
}
