package repository

import (
	"auth/internal/domain/entity"
	"context"
	"database/sql"
	"errors"
)

type TokenRepository struct {
	DB *sql.DB
}

func NewTokenRepository(db *sql.DB) *TokenRepository {
	return &TokenRepository{
		DB: db,
	}
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
        SELECT id, user_id, access_token, refresh_token, 
               access_token_expires_at, refresh_token_expires_at
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

// GetTokensByUserID retrieves all tokens for a user
func (r *TokenRepository) GetTokensByUserID(ctx context.Context, tx *sql.Tx, userID int) ([]*entity.RefreshToken, error) {
	query := `
        SELECT id, user_id, access_token, refresh_token, 
               access_token_expires_at, refresh_token_expires_at
        FROM auth_tokens
        WHERE user_id = $1
    `

	rows, err := tx.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tokens []*entity.RefreshToken
	for rows.Next() {
		token := &entity.RefreshToken{}
		err := rows.Scan(
			&token.ID,
			&token.UserID,
			&token.RefreshToken,
			&token.RefreshTokenExpiredAt,
		)
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, token)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return tokens, nil
}

// DeleteToken deletes a token by its ID
func (r *TokenRepository) DeleteToken(ctx context.Context, tx *sql.Tx, tokenID int) error {
	query := `
        DELETE FROM auth_tokens
        WHERE id = $1
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
		return errors.New("no rows affected, token might not exist")
	}

	return nil
}
