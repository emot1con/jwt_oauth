package domain_interface

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
	Delete(ctx context.Context, tx *sql.Tx, id string) error
}

type UserServiceInterface interface {
	Create(ctx context.Context, tx *sql.Tx, user *entity.User) error
	GetByID(ctx context.Context, tx *sql.Tx, id int) (*entity.User, error)
	GetByEmail(ctx context.Context, tx *sql.Tx, email string) (*entity.User, error)
	Update(ctx context.Context, tx *sql.Tx, user *entity.User) error
	Delete(ctx context.Context, tx *sql.Tx, id string) error
}

type UserUsecaseInterface interface {
	Register(payload *entity.RegisterPayload) (*entity.RegisterResponse, error)
	Login(payload *entity.LoginPayload) (*entity.JWTResponse, error)
	GetByID(payload int) (entity.User, error)
	GetByEmail(payload string) (*entity.User, error)
	Logout(ID int) error
}
