package repository

import (
	"auth/internal/domain/entity"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/sirupsen/logrus"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository struct{}

func NewUserRepository() *UserRepository { return &UserRepository{} }

// Create inserts a new user into the database
func (r *UserRepository) Create(ctx context.Context, tx *sql.Tx, user *entity.User) error {
	logrus.Info("inserting new user to database repository with user: ", user)
	query := `
		INSERT INTO users (email, password, name, provider, provider_id)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := tx.ExecContext(ctx, query,
		user.Email,
		user.Password,
		user.Name,
		user.Provider,
		user.ProviderID,
	)

	return err
}

// GetByID retrieves a user by their ID
func (r *UserRepository) GetByID(ctx context.Context, tx *sql.Tx, id int) (*entity.User, error) {
	query := `
		SELECT id, email, password, name, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &entity.User{}
	err := tx.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

// GetByEmail retrieves a user by their email
func (r *UserRepository) GetByEmail(ctx context.Context, tx *sql.Tx, email string) (*entity.User, error) {
	logrus.Info("repo facebook email: ", email)
	query := `
		SELECT id, email, password, name, created_at, updated_at
		FROM users
		WHERE email = $1
	`

	user := &entity.User{}
	logrus.Info("repo facebook email query: ", email)

	err := tx.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Password,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		logrus.Infof("failed repo query: %v with error: %v", email, err)
		return nil, err
	}
	logrus.Info("success repo query: ", email)

	return user, nil
}

// Update modifies an existing user in the database
func (r *UserRepository) Update(ctx context.Context, tx *sql.Tx, user *entity.User) error {
	query := `
		UPDATE users
		SET email = $1, password = $2, name = $3, provider = $4, provider_id = $5, updated_at = $6
		WHERE id = $7
	`

	_, err := tx.ExecContext(ctx, query,
		user.Email,
		user.Password,
		user.Name,
		user.Provider,
		user.ProviderID,
		time.Now(),
		user.ID,
	)

	return err
}

// Delete removes a user from the database
func (r *UserRepository) Delete(ctx context.Context, tx *sql.Tx, id int) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	result, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}
