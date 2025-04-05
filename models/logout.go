package models

import (
	"database/sql"
	"errors"
	"time"
)

var (
	ErrLogoutAlreadyExists = errors.New("logout already exists")
	ErrLogoutNotFound      = errors.New("logout not found")
)

type Logout struct {
	ID        string
	UserID    sql.NullString
	Token     string
	RevokedAt time.Time
}
