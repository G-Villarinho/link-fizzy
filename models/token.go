package models

type TokenClaims struct {
	Sub string `json:"sub"`
	Sid string `json:"sid"`
}
