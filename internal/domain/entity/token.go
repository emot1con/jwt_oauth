package entity

type RefreshToken struct {
	ID                    int    `json:"id"`
	UserID                int    `json:"user_id"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiredAt string `json:"refresh_token_expired_at"`
}
