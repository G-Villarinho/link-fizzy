package services

import (
	"context"
	"fmt"

	"github.com/g-villarinho/link-fizz-api/models"
	"github.com/g-villarinho/link-fizz-api/pkgs/di"
)

type AuthService interface {
	Login(ctx context.Context, email, password string) (*models.LoginResponse, error)
	Logout(ctx context.Context, userID, token string) error
}

type authService struct {
	i  *di.Injector
	ss SecurityService
	ts TokenService
	ls LogoutService
	us UserService
}

func NewAuthService(i *di.Injector) (AuthService, error) {
	securityService, err := di.Invoke[SecurityService](i)
	if err != nil {
		return nil, fmt.Errorf("invoke services.SecurityService: %w", err)
	}

	tokenService, err := di.Invoke[TokenService](i)
	if err != nil {
		return nil, fmt.Errorf("invoke services.TokenService: %w", err)
	}

	logoutService, err := di.Invoke[LogoutService](i)
	if err != nil {
		return nil, fmt.Errorf("invoke services.LogoutService: %w", err)
	}

	userService, err := di.Invoke[UserService](i)
	if err != nil {
		return nil, fmt.Errorf("invoke services.UserService: %w", err)
	}

	return &authService{
		i:  i,
		ss: securityService,
		ts: tokenService,
		ls: logoutService,
		us: userService,
	}, nil
}

func (l *authService) Login(ctx context.Context, email string, password string) (*models.LoginResponse, error) {
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

func (l *authService) Logout(ctx context.Context, userID string, token string) error {
	logout, err := l.ls.IsLogoutRegistered(ctx, token)
	if err != nil {
		return fmt.Errorf("check logout: %w", err)
	}

	if logout {
		return models.ErrLogoutAlreadyExists
	}

	if err := l.ls.CreateLogout(ctx, &userID, token); err != nil {
		return fmt.Errorf("create logout: %w", err)
	}

	return nil
}
