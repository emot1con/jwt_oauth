package services

import (
	"auth/internal/domain/entity"
	"auth/internal/domain/interface"
	"context"
	"database/sql"
)

type TokenService struct {
	repo interfaces.TokenRepsitoryInterface
	Tx   *sql.Tx
	ctx  context.Context
}

func NewTokenService(repo interfaces.TokenRepsitoryInterface, DB *sql.Tx, ctx context.Context) *TokenService {
	return &TokenService{
		repo: repo,
		Tx:   DB,
		ctx:  ctx,
	}
}

func (s *TokenService) SaveToken(refreshToken *entity.RefreshToken) error {
	return s.repo.SaveToken(s.ctx, s.Tx, refreshToken)
}

func (s *TokenService) GetTokenByRefresh(refreshToken string) (*entity.RefreshToken, error) {
	return s.repo.GetTokenByRefresh(s.ctx, s.Tx, refreshToken)
}

func (s *TokenService) GetTokensByUserID(userID int) ([]*entity.RefreshToken, error) {
	return s.repo.GetTokensByUserID(s.ctx, s.Tx, userID)
}

func (s *TokenService) DeleteToken(tokenID int) error {
	return s.repo.DeleteToken(s.ctx, s.Tx, tokenID)
}
