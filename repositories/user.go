package repositories

import (
	"context"
	"database/sql"

	"github.com/g-villarinho/link-fizz-api/models"
	"github.com/g-villarinho/link-fizz-api/pkgs/di"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error
	DeleteUser(ctx context.Context, id string) error
}

type userRepository struct {
	i  *di.Injector
	db *sql.DB
}

func NewUserRepository(i *di.Injector) (UserRepository, error) {
	db, err := di.Invoke[*sql.DB](i)
	if err != nil {
		return nil, err
	}

	return &userRepository{
		i:  i,
		db: db,
	}, nil
}

func (u *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	statement, err := u.db.PrepareContext(ctx, "INSERT INTO users (id, name, email, password_hash, created_at) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.ExecContext(ctx, user.ID, user.Name, user.Email, user.PasswordHash, user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	statement, err := u.db.PrepareContext(ctx, "SELECT id, name, email, password_hash, created_at FROM users WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	return scanUser(statement.QueryRowContext(ctx, id))
}

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	statement, err := u.db.PrepareContext(ctx, "SELECT id, name, email, password_hash, created_at FROM users WHERE email = ?")
	if err != nil {
		return nil, err
	}
	defer statement.Close()

	return scanUser(statement.QueryRowContext(ctx, email))
}

func (u *userRepository) UpdateUser(ctx context.Context, user *models.User) error {
	statement, err := u.db.PrepareContext(ctx, "UPDATE users SET name = ?, email = ?, password_hash = ?, updated_at = ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.ExecContext(ctx, user.Name, user.Email, user.PasswordHash, user.UpdatedAt, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (u *userRepository) DeleteUser(ctx context.Context, id string) error {
	statement, err := u.db.PrepareContext(ctx, "DELETE FROM users WHERE id = ?")
	if err != nil {
		return err
	}
	defer statement.Close()

	_, err = statement.ExecContext(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func scanUser(row *sql.Row) (*models.User, error) {
	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}
