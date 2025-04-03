package models

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type LoginPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
