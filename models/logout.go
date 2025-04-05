package models

import (
	"errors"
	"time"
)

var (
	ErrLogoutAlreadyExists = errors.New("logout already exists")
	ErrLogoutNotFound      = errors.New("logout not found")
)

type Logout struct {
	ID        string
	UserID    string
	Token     string
	RevokedAt time.Time
}
