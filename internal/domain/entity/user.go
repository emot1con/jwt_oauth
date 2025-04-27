package entity

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RegisterPayload struct {
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type RegisterResponse struct {
	Message string `json:"message"`
}

type JWTResponse struct {
	Token                 string `json:"token"`
	TokenExpiredAt        string `json:"token_expired_at"`
	RefreshToken          string `json:"refresh_token"`
	RefreshTokenExpiredAt string `json:"refresh_token_expired_at"`
}
