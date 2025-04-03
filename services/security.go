package services

import (
	"context"

	"github.com/g-villarinho/link-fizz-api/pkgs/di"
	"golang.org/x/crypto/bcrypt"
)

type SecurityService interface {
	HashPassword(ctx context.Context, password string) (string, error)
	VerifyPassword(ctx context.Context, hashedPassword, password string) error
}

type securityService struct {
	i *di.Injector
}

func NewSecurityService(i *di.Injector) (SecurityService, error) {
	return &securityService{
		i: i,
	}, nil
}

func (s *securityService) HashPassword(ctx context.Context, password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (s *securityService) VerifyPassword(ctx context.Context, hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
