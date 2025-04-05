package services

import (
	"context"
	"fmt"
	"time"

	"github.com/g-villarinho/link-fizz-api/models"
	"github.com/g-villarinho/link-fizz-api/pkgs/di"
	"github.com/g-villarinho/link-fizz-api/repositories"
	"github.com/google/uuid"
)

type UserService interface {
	CreateUser(ctx context.Context, name, email, password string) error
	GetUserByID(ctx context.Context, ID string) (*models.UserResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, ID string, name, email string) error
	DeleteUser(ctx context.Context, ID, token string) error
}

type userService struct {
	i  *di.Injector
	ss SecurityService
	ls LogoutService
	ur repositories.UserRepository
}

func NewUserService(i *di.Injector) (UserService, error) {
	securityService, err := di.Invoke[SecurityService](i)
	if err != nil {
		return nil, fmt.Errorf("invoke services.security: %w", err)
	}

	logoutService, err := di.Invoke[LogoutService](i)
	if err != nil {
		return nil, fmt.Errorf("invoke services.logout: %w", err)
	}

	userRepository, err := di.Invoke[repositories.UserRepository](i)
	if err != nil {
		return nil, fmt.Errorf("invoke repositories.user: %w", err)
	}

	return &userService{
		i:  i,
		ss: securityService,
		ls: logoutService,
		ur: userRepository,
	}, nil
}

func (u *userService) CreateUser(ctx context.Context, name, email, password string) error {
	userFromEmail, err := u.ur.GetUserByEmail(ctx, email)
	if err != nil {
		return fmt.Errorf("get user by email: %w", err)
	}

	if userFromEmail != nil {
		return models.ErrUserAlreadyExists
	}

	passwordHash, err := u.ss.HashPassword(ctx, password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	id, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("generate UUID: %w", err)
	}

	user := models.NewUser(id.String(), name, email, passwordHash, time.Now().UTC())

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

func (u *userService) UpdateUser(ctx context.Context, ID string, name string, email string) error {
	user, err := u.ur.GetUserByID(ctx, ID)
	if err != nil {
		return fmt.Errorf("get user by ID %s: %w", ID, err)
	}

	if user == nil {
		return models.ErrUserNotFound
	}

	if user.Email != email {
		userFromEmail, err := u.ur.GetUserByEmail(ctx, email)
		if err != nil {
			return fmt.Errorf("get user by email: %w", err)
		}

		if userFromEmail != nil {
			return models.ErrUserAlreadyExists
		}

		user.Email = email
	}

	user.Name = name

	if err := u.ur.UpdateUser(ctx, user); err != nil {
		return fmt.Errorf("update user: %w", err)
	}

	return nil
}

func (u *userService) DeleteUser(ctx context.Context, ID, token string) error {
	if err := u.ur.DeleteUser(ctx, ID); err != nil {
		return fmt.Errorf("delete user: %w", err)
	}

	if err := u.ls.CreateLogout(ctx, nil, token); err != nil {
		return fmt.Errorf("create logout: %w", err)
	}

	return nil
}
