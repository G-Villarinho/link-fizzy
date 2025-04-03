package models

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrLinkNotFound     = errors.New("link not found")
	ErrInvalidShortCode = errors.New("invalid short code")
)

type Link struct {
	ID          string
	OriginalURL string
	ShortCode   string
	UserID      string
	CreatedAt   time.Time
	UpdatedAt   sql.NullTime
}

type LinkPayload struct {
	OriginalURL string `json:"original_url"`
}
