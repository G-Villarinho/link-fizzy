package services

import (
	"context"
	"fmt"

	"github.com/g-villarinho/link-fizz-api/models"
	"github.com/g-villarinho/link-fizz-api/pkgs/di"
)

type LoginService interface {
	Login(ctx context.Context, email, password string) (*models.LoginResponse, error)
}

type loginService struct {
	i  *di.Injector
	ss SecurityService
	ts TokenService
	us UserService
}

func NewLoginService(i *di.Injector) (LoginService, error) {
	securityService, err := di.Invoke[SecurityService](i)
	if err != nil {
		return nil, fmt.Errorf("invoke services.SecurityService: %w", err)
	}

	tokenService, err := di.Invoke[TokenService](i)
	if err != nil {
		return nil, fmt.Errorf("invoke services.TokenService: %w", err)
	}

	userService, err := di.Invoke[UserService](i)
	if err != nil {
		return nil, fmt.Errorf("invoke services.UserService: %w", err)
	}

	return &loginService{
		i:  i,
		ss: securityService,
		ts: tokenService,
		us: userService,
	}, nil
}

func (l *loginService) Login(ctx context.Context, email string, password string) (*models.LoginResponse, error) {
	user, err := l.us.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := l.ss.VerifyPassword(ctx, user.PasswordHash, password); err != nil {
		return nil, models.ErrInvalidCredentials
	}

	token, err := l.ts.GenerateToken(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("generate token: %w", err)
	}

	return &models.LoginResponse{
		Token: token,
	}, nil
}
