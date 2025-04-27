package usecases

import (
	"auth/internal/config"
	"auth/internal/domain/entity"
	domain_interface "auth/internal/domain/interface"
	"auth/pkg/helper"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
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

	accessToken, err := helper.GenerateToken(user.ID, "user", 1, 0, 0)
	if err != nil {
		return nil, err
	}

	refreshToken, err := helper.GenerateToken(user.ID, "user", 0, 3, 0)
	if err != nil {
		return nil, err
	}

	return &entity.JWTResponse{
		Token:                 fmt.Sprintf("Bearer %s", accessToken),
		TokenExpiredAt:        time.Now().AddDate(1, 0, 0).Format("2006-01-02 15:04:05"),
		RefreshToken:          fmt.Sprintf("Bearer %s", refreshToken),
		RefreshTokenExpiredAt: time.Now().AddDate(0, 3, 0).Format("2006-01-02 15:04:05"),
	}, nil
}

func (s *UserUseCase) GetUserByID(id int) (*entity.User, error) {
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

	ctx := context.Background()
	user, err := s.userService.GetByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserUseCase) DeleteUser(payload *entity.User) error {
	tx, err := s.dB.Begin()
	if err != nil {
		return err
	}

	defer helper.CommitOrRollback(tx)
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	ctx := context.Background()
	err = s.userService.Delete(ctx, tx, payload.ID)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserUseCase) Logout(bearerToken string) error {
	tx, err := s.dB.Begin()
	if err != nil {
		return err
	}

	defer helper.CommitOrRollback(tx)
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	stringToken := strings.Split(bearerToken, " ")[1]
	token, err := jwt.ParseWithClaims(stringToken, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil || !token.Valid {
		return errors.New("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("invalid claims")
	}

	expAt := int64(claims["exp"].(float64))
	if time.Now().Unix() > expAt {
		return fmt.Errorf("acces token expired")
	}

	logrus.Info("blacklist token")
	ttl := time.Until(time.Unix(expAt, 0))
	redisKey := "blacklist:" + strings.Split(stringToken, " ")[1]
	if err := config.RedisClient.Set(context.Background(), redisKey, "revoked", ttl).Err(); err != nil {
		return err
	}

	return nil
}
