package interfaces

import (
	"auth/internal/domain/entity"
	"context"
	"database/sql"
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, tx *sql.Tx, user *entity.User) error
	GetByID(ctx context.Context, tx *sql.Tx, id int) (*entity.User, error)
	GetByEmail(ctx context.Context, tx *sql.Tx, email string) (*entity.User, error)
	Update(ctx context.Context, tx *sql.Tx, user *entity.User) error
	Delete(ctx context.Context, tx *sql.Tx, id int) error
}

type UserServiceInterface interface {
	Create(ctx context.Context, tx *sql.Tx, user *entity.User) error
	GetByID(ctx context.Context, tx *sql.Tx, id int) (*entity.User, error)
	GetByEmail(ctx context.Context, tx *sql.Tx, email string) (*entity.User, error)
	Update(ctx context.Context, tx *sql.Tx, user *entity.User) error
	Delete(ctx context.Context, tx *sql.Tx, id int) error
	RefreshToken(ctx context.Context, refreshToken string) (string, error)
}

type UserUsecaseInterface interface {
	Register(payload *entity.RegisterPayload) (*entity.RegisterResponse, error)
	Login(payload *entity.LoginPayload) (*entity.JWTResponse, error)
	GetUserByID(payload int) (entity.User, error)
	Logout(bearerToken string) error
	DeleteUser(payload int) error
}

type TokenRepsitoryInterface interface {
	SaveToken(ctx context.Context, tx *sql.Tx, refreshToken *entity.RefreshToken) error
	GetTokenByRefresh(ctx context.Context, tx *sql.Tx, refreshToken string) (*entity.RefreshToken, error)
	GetTokensByUserID(ctx context.Context, tx *sql.Tx, userID int) ([]*entity.RefreshToken, error)
	DeleteToken(ctx context.Context, tx *sql.Tx, tokenID int) error
}
