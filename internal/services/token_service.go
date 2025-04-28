package services

import (
	"auth/internal/domain/entity"
	"auth/internal/domain/interface"
	"context"
	"database/sql"
)

type TokenService struct {
	repo interfaces.TokenRepositoryInterface
}

func NewTokenService(repo interfaces.TokenRepositoryInterface) *TokenService {
	return &TokenService{
		repo: repo,
	}
}

func (s *TokenService) SaveToken(refreshToken *entity.RefreshToken, Tx *sql.Tx, ctx context.Context) error {
	return s.repo.SaveToken(ctx, Tx, refreshToken)
}

func (s *TokenService) GetTokenByRefresh(refreshToken string, Tx *sql.Tx, ctx context.Context) (*entity.RefreshToken, error) {
	return s.repo.GetTokenByRefresh(ctx, Tx, refreshToken)
}

func (s *TokenService) GetTokensByUserID(userID int, Tx *sql.Tx, ctx context.Context) ([]*entity.RefreshToken, error) {
	return s.repo.GetTokensByUserID(ctx, Tx, userID)
}

func (s *TokenService) DeleteToken(tokenID int, Tx *sql.Tx, ctx context.Context) error {
	return s.repo.DeleteToken(ctx, Tx, tokenID)
}
