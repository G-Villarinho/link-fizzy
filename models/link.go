package models

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrLinkNotFound            = errors.New("link not found")
	ErrInvalidShortCode        = errors.New("invalid short code")
	ErrCustomCodeAlreadyExists = errors.New("custom code already exists")
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
	DestinationURL string  `json:"destination_url"`
	CustomCode     *string `json:"custom_code,omitempty"`
}
