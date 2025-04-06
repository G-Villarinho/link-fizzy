package models

import (
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

var (
	ErrLinkNotFound            = errors.New("link not found")
	ErrInvalidShortCode        = errors.New("invalid short code")
	ErrCustomCodeAlreadyExists = errors.New("custom code already exists")
	ErrLinkNotBelongToUser     = errors.New("link does not belong to user")
)

type Link struct {
	ID          string
	Title       sql.NullString
	OriginalURL string
	ShortCode   string
	UserID      string
	CreatedAt   time.Time
	UpdatedAt   sql.NullTime
}

type LinkPayload struct {
	Title          *string `json:"title,omitempty"`
	DestinationURL string  `json:"destinationUrl"`
	CustomCode     *string `json:"customCode,omitempty"`
}

type CreateLinkResponse struct {
	ShortCode string `json:"shortCode"`
}

type LinkResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	OriginalURL string `json:"originalUrl"`
	ShortCode   string `json:"shortCode"`
	ShortURL    string `json:"shortUrl"`
	CreatedAt   string `json:"createdAt"`
}

func (l *Link) ToResponse(apiURL string) LinkResponse {
	var title string

	if l.Title.Valid && strings.TrimSpace(l.Title.String) != "" {
		title = l.Title.String
	} else {
		parsedURL, err := url.Parse(l.OriginalURL)
		if err != nil {
			title = "untitled"
		} else {
			domain := parsedURL.Hostname()
			domain = strings.TrimPrefix(domain, "www.")
			title = domain + " - untitled"
		}
	}

	return LinkResponse{
		ID:          l.ID,
		Title:       title,
		OriginalURL: l.OriginalURL,
		ShortCode:   l.ShortCode,
		ShortURL:    fmt.Sprintf("%s/%s", apiURL, l.ShortCode),
		CreatedAt:   l.CreatedAt.Format(time.RFC3339),
	}
}
