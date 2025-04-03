package models

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUserAlreadyExists = fmt.Errorf("user already exists")
	ErrUserNotFound      = fmt.Errorf("user not found")
)

type User struct {
	ID           string
	Name         string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    sql.NullString
}

type CreateUserPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (p *CreateUserPayload) ToUser() (*User, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return nil, fmt.Errorf("generate UUID: %w", err)
	}

	return &User{
		ID:        id.String(),
		Name:      p.Name,
		Email:     p.Email,
		CreatedAt: time.Now().UTC(),
	}, nil
}

func (u *User) ToUseResponse() *UserResponse {
	return &UserResponse{
		ID:    u.ID,
		Name:  u.Name,
		Email: u.Email,
	}
}
