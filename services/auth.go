package services

import (
	"context"
	"fmt"

	"github.com/g-villarinho/link-fizz-api/models"
	"github.com/g-villarinho/link-fizz-api/pkgs/di"
)

type AuthService interface {
	Register(ctx context.Context, name, email, password, ipAdress, userAgent string) (*models.LoginResponse, error)
	Login(ctx context.Context, email, password, ipAdress, userAgent string) (*models.LoginResponse, error)
	Logout(ctx context.Context, userID, sessionID, token string) error
}

type authService struct {
	i   *di.Injector
	scs SecurityService
	ss  SessionService
	ls  LogoutService
	us  UserService
}

func NewAuthService(i *di.Injector) (AuthService, error) {
	securityService, err := di.Invoke[SecurityService](i)
	if err != nil {
		return nil, fmt.Errorf("invoke services.SecurityService: %w", err)
	}

	sessionService, err := di.Invoke[SessionService](i)
	if err != nil {
		return nil, fmt.Errorf("invoke services.SessionService: %w", err)
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
		i:   i,
		scs: securityService,
		ss:  sessionService,
		ls:  logoutService,
		us:  userService,
	}, nil
}

func (l *authService) Register(ctx context.Context, name, email, password, ipAdress, userAgent string) (*models.LoginResponse, error) {
	userID, err := l.us.CreateUser(ctx, name, email, password)
	if err != nil {
		return nil, err
	}

	token, err := l.ss.CreateSession(ctx, userID, ipAdress, userAgent)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	return &models.LoginResponse{
		Token: token,
	}, nil
}

func (l *authService) Login(ctx context.Context, email, password, ipAdress, userAgent string) (*models.LoginResponse, error) {
	user, err := l.us.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}

	if err := l.scs.VerifyPassword(ctx, user.PasswordHash, password); err != nil {
		return nil, models.ErrInvalidCredentials
	}

	token, err := l.ss.CreateSession(ctx, user.ID, ipAdress, userAgent)
	if err != nil {
		return nil, fmt.Errorf("create session: %w", err)
	}

	return &models.LoginResponse{
		Token: token,
	}, nil
}

func (l *authService) Logout(ctx context.Context, userID, sessionID, token string) error {
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

	if err := l.ss.DeleteSession(ctx, sessionID); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}
