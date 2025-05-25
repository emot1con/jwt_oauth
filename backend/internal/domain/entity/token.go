package entity

import "time"

type RefreshToken struct {
	ID                    int       `json:"id"`
	UserID                int       `json:"user_id"`
	RefreshToken          string    `json:"refresh_token" validate:"required"`
	RefreshTokenExpiredAt time.Time `json:"refresh_token_expired_at" validate:"required"`
}
