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
	ID          string       `json:"id"`
	OriginalURL string       `json:"original_url"`
	ShortCode   string       `json:"short_code"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   sql.NullTime `json:"updated_at"`
}

type LinkPayload struct {
	OriginalURL string `json:"original_url"`
}
