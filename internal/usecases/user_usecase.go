package usecases

import (
	"auth/internal/domain/entity"
	"auth/internal/domain/interface"
	"auth/pkg/helper"
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type UserUseCase struct {
	userService domain_interface.UserServiceInterface
	dB          *sql.DB
	validator   *validator.Validate
}

func NewUserUseCase(userService domain_interface.UserServiceInterface, DB *sql.DB) *UserUseCase {
	return &UserUseCase{
		userService: userService,
		dB:          DB,
	}
}

func (s *UserUseCase) Register(payload *entity.RegisterPayload) (*entity.RegisterResponse, error) {
	tx, err := s.dB.Begin()
	if err != nil {
		return nil, err
	}

	defer helper.CommitOrRollback(tx)
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	if err := s.validator.Struct(payload); err != nil {
		return nil, err
	}

	ctx := context.Background()
	if _, err := s.userService.GetByEmail(ctx, tx, payload.Email); err == nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := helper.GenerateHashPassword(payload.Password)
	if err != nil {
		return nil, err
	}

	newUser := &entity.User{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: string(hashedPassword),
	}

	err = s.userService.Create(ctx, tx, newUser)
	if err != nil {
		return nil, err
	}

	return &entity.RegisterResponse{
		Message: "User registered successfully",
	}, nil
}

func (s *UserUseCase) Login(payload *entity.LoginPayload) (*entity.JWTResponse, error) {
	tx, err := s.dB.Begin()
	if err != nil {
		return nil, err
	}

	defer helper.CommitOrRollback(tx)
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	logrus.Info("get user by email")
	ctx := context.Background()
	user, err := s.userService.GetByEmail(ctx, tx, payload.Email)
	if err != nil {
		return nil, err
	}

	logrus.Info("compare password")
	if !helper.ComparePassword(user.Password, []byte(payload.Password)) {
		return nil, errors.New("invalid password")
	}

	token, err := helper.GenerateToken(user.ID)
	if err != nil {
		return nil, err
	}

	return &entity.JWTResponse{
		Token:     token,
		ExpiredAt: time.Now().Add(time.Hour * 24),
	}, nil

}
