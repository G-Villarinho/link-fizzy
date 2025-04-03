package services

import (
	"context"
	"fmt"

	"github.com/g-villarinho/link-fizz-api/models"
	"github.com/g-villarinho/link-fizz-api/pkgs/di"
	"github.com/g-villarinho/link-fizz-api/repositories"
)

type UserService interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, ID string) (*models.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type userService struct {
	i  *di.Injector
	ss SecurityService
	ur repositories.UserRepository
}

func NewUserService(i *di.Injector) (UserService, error) {
	securityService, err := di.Invoke[SecurityService](i)
	if err != nil {
		return nil, err
	}

	userRepository, err := di.Invoke[repositories.UserRepository](i)
	if err != nil {
		return nil, err
	}

	return &userService{
		i:  i,
		ss: securityService,
		ur: userRepository,
	}, nil
}

func (u *userService) CreateUser(ctx context.Context, user *models.User) error {
	userFromEmail, err := u.ur.GetUserByEmail(ctx, user.Email)
	if err != nil {
		return fmt.Errorf("get user by email: %w", err)
	}

	if userFromEmail != nil {
		return models.ErrUserAlreadyExists
	}

	passwordHash, err := u.ss.HashPassword(ctx, user.PasswordHash)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	user.PasswordHash = passwordHash

	if err := u.ur.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (u *userService) GetUserByID(ctx context.Context, ID string) (*models.UserResponse, error) {
	user, err := u.ur.GetUserByID(ctx, ID)
	if err != nil {
		return nil, fmt.Errorf("get user by ID %s: %w", ID, err)
	}

	if user == nil {
		return nil, models.ErrUserNotFound
	}

	return user.ToUseResponse(), nil
}

func (u *userService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	user, err := u.ur.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("get user by email %s: %w", email, err)
	}

	if user == nil {
		return nil, models.ErrUserNotFound
	}

	return user, nil
}
