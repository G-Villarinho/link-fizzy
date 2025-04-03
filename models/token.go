package models

import (
	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	Sub string `json:"sub"`
	jwt.RegisteredClaims
}
