package repository

import (
	"auth/internal/domain/entity"
	"context"
	"database/sql"
	"errors"
)

type TokenRepository struct{}

func NewTokenRepository() *TokenRepository {
	return &TokenRepository{}
}

// SaveToken saves a new token to the database
func (r *TokenRepository) SaveToken(ctx context.Context, tx *sql.Tx, refreshToken *entity.RefreshToken) error {
	query := `
        INSERT INTO auth_tokens (
            user_id, 
            refresh_token, 
            refresh_token_expires_at
        ) VALUES ($1, $2, $3)
    `

	_, err := tx.ExecContext(
		ctx,
		query,
		refreshToken.UserID,
		refreshToken.RefreshToken,
		refreshToken.RefreshTokenExpiredAt,
	)

	return err
}

// GetTokenByRefresh retrieves a token by refresh token
func (r *TokenRepository) GetTokenByRefresh(ctx context.Context, tx *sql.Tx, refreshToken string) (*entity.RefreshToken, error) {
	query := `
        SELECT id, user_id, refresh_token, 
                refresh_token_expires_at
        FROM auth_tokens
        WHERE refresh_token = $1
    `

	token := &entity.RefreshToken{}
	err := tx.QueryRowContext(ctx, query, refreshToken).Scan(
		&token.ID,
		&token.UserID,
		&token.RefreshToken,
		&token.RefreshTokenExpiredAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors.New("token not found")
		}
		return nil, err
	}

	return token, nil
}

// GetTokenByUserID retrieves the most recent token for a user
func (r *TokenRepository) GetTokenByUserID(ctx context.Context, tx *sql.Tx, userID int) (*entity.RefreshToken, error) {
	query := `
        SELECT id, user_id, refresh_token, refresh_token_expires_at
        FROM auth_tokens
        WHERE user_id = $1
        ORDER BY refresh_token_expires_at DESC
        LIMIT 1
    `

	token := &entity.RefreshToken{}
	err := tx.QueryRowContext(ctx, query, userID).Scan(
		&token.ID,
		&token.UserID,
		&token.RefreshToken,
		&token.RefreshTokenExpiredAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Return nil instead of error when no token is found
		}
		return nil, err
	}

	return token, nil
}

func (r *TokenRepository) UpdateToken(ctx context.Context, tx *sql.Tx, token *entity.RefreshToken) error {
	query := `
        UPDATE auth_tokens 
        SET refresh_token = $1, 
            refresh_token_expires_at = $2
        WHERE id = $3
    `

	result, err := tx.ExecContext(
		ctx,
		query,
		token.RefreshToken,
		token.RefreshTokenExpiredAt,
		token.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("token not found")
	}

	return nil
}

// DeleteToken deletes a token by its ID
func (r *TokenRepository) DeleteToken(ctx context.Context, tx *sql.Tx, tokenID int) error {
	query := `
        DELETE FROM auth_tokens
        WHERE user_id = $1
    `

	result, err := tx.ExecContext(ctx, query, tokenID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("invalid token")
	}

	return nil
}
