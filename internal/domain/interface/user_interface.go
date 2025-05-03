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
	Create(ctx context.Context, tx *sql.Tx, payload *entity.RegisterPayload) error
	GetByID(ctx context.Context, tx *sql.Tx, id int) (*entity.User, error)
	GetByEmail(ctx context.Context, tx *sql.Tx, email string) (*entity.User, error)
	Update(ctx context.Context, tx *sql.Tx, user *entity.User) error
	Delete(ctx context.Context, tx *sql.Tx, id int) error
}

type UserUsecaseInterface interface {
	Register(payload *entity.RegisterPayload) (*entity.RegisterResponse, error)
	Login(payload *entity.LoginPayload) (*entity.JWTResponse, error)
	GetUserByID(ID int) (*entity.User, error)
	Logout(ID int, token string) error
	RefreshToken(refreshToken string) (*entity.JWTResponse, error)
	DeleteUser(ID int) error
}

type TokenRepositoryInterface interface {
	SaveToken(ctx context.Context, tx *sql.Tx, refreshToken *entity.RefreshToken) error
	GetTokenByRefresh(ctx context.Context, tx *sql.Tx, refreshToken string) (*entity.RefreshToken, error)
	GetTokenByUserID(ctx context.Context, tx *sql.Tx, userID int) (*entity.RefreshToken, error)
	UpdateToken(ctx context.Context, tx *sql.Tx, token *entity.RefreshToken) error
	DeleteToken(ctx context.Context, tx *sql.Tx, tokenID int) error
}

type TokenServiceInterface interface {
	SaveToken(refreshToken *entity.RefreshToken, Tx *sql.Tx, ctx context.Context) error
	GetTokenByRefresh(refreshToken string, Tx *sql.Tx, ctx context.Context) (*entity.RefreshToken, error)
	GetTokenByUserID(userID int, Tx *sql.Tx, ctx context.Context) (*entity.RefreshToken, error)
	UpdateToken(token *entity.RefreshToken, Tx *sql.Tx, ctx context.Context) error
	DeleteToken(tokenID int, Tx *sql.Tx, ctx context.Context) error
}
