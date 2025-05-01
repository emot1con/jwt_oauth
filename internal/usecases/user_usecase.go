package usecases

import (
	"auth/internal/domain/entity"
	"auth/internal/domain/interface"
	"auth/internal/repository"
	"auth/pkg/helper"
	"auth/pkg/middleware"
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
	UserService  interfaces.UserServiceInterface
	TokenService interfaces.TokenServiceInterface
	DB           *sql.DB
	validator    *validator.Validate
}

func NewUserUseCase(
	userService interfaces.UserServiceInterface,
	tokenService interfaces.TokenServiceInterface,
	DB *sql.DB,
	validator *validator.Validate,
) *UserUseCase {
	return &UserUseCase{
		UserService:  userService,
		DB:           DB,
		TokenService: tokenService,
		validator:    validator,
	}
}

func (s *UserUseCase) Register(payload *entity.RegisterPayload) (*entity.RegisterResponse, error) {
	ctx := context.Background()

	var registerResp *entity.RegisterResponse
	if err := middleware.WithTransaction(ctx, s.DB, func(tx *sql.Tx) error {
		if err := s.validator.Struct(payload); err != nil {
			return err
		}

		ctx := context.Background()

		logrus.Info("checking if email already exists")
		_, err := s.UserService.GetByEmail(ctx, tx, payload.Email)
		if err == nil {
			return errors.New("email already exists")
		} else if err != sql.ErrNoRows && err != repository.ErrUserNotFound {
			return err
		}

		hashedPassword, err := helper.GenerateHashPassword(payload.Password)
		if err != nil {
			return err
		}
		payload.Password = string(hashedPassword)

		logrus.Info("insert new user to database")
		if err := s.UserService.Create(ctx, tx, payload); err != nil {
			return err
		}

		registerResp = &entity.RegisterResponse{
			Message: "User registered successfully",
		}
		return nil
	}); err != nil {
		return nil, err
	}

	return registerResp, nil
}

func (s *UserUseCase) Login(payload *entity.LoginPayload) (*entity.JWTResponse, error) {
	ctx := context.Background()

	var jwtResp *entity.JWTResponse
	if err := middleware.WithTransaction(ctx, s.DB, func(tx *sql.Tx) error {
		logrus.Info("get user by email")
		user, err := s.UserService.GetByEmail(ctx, tx, payload.Email)
		if err != nil {
			return fmt.Errorf("email not found")
		}

		if err := helper.ComparePassword(user.Password, []byte(payload.Password)); err != nil {
			return errors.New("wrong password")
		}

		accessToken, err := helper.GenerateToken(user.ID, "user", 1, 0, 0)
		if err != nil {
			return err
		}

		refreshToken, err := helper.GenerateToken(user.ID, "user", 0, 3, 0)
		if err != nil {
			return err
		}

		logrus.Info("save refresh token to database")
		err = s.TokenService.SaveToken(&entity.RefreshToken{
			UserID:                user.ID,
			RefreshToken:          refreshToken,
			RefreshTokenExpiredAt: time.Now().AddDate(0, 3, 0),
		}, tx, ctx)
		if err != nil {
			return err
		}
		logrus.Info("refresh token created successfully and saved to database")

		jwtResp = &entity.JWTResponse{
			Token:                 fmt.Sprintf("Bearer %s", accessToken),
			TokenExpiredAt:        time.Now().AddDate(1, 0, 0).Format("2006-01-02 15:04:05"),
			RefreshToken:          fmt.Sprintf("Bearer %s", refreshToken),
			RefreshTokenExpiredAt: time.Now().AddDate(0, 3, 0).Format("2006-01-02 15:04:05"),
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return jwtResp, nil
}

func (s *UserUseCase) GetUserByID(id int) (*entity.User, error) {
	ctx := context.Background()

	var userResp *entity.User
	if err := middleware.WithTransaction(ctx, s.DB, func(tx *sql.Tx) error {
		logrus.Info("find user by id")
		user, err := s.UserService.GetByID(ctx, tx, id)
		if err != nil {
			return err
		}
		userResp = user

		return nil
	}); err != nil {
		return nil, err
	}

	return userResp, nil
}

func (s *UserUseCase) DeleteUser(ID int) error {
	ctx := context.Background()

	if err := middleware.WithTransaction(ctx, s.DB, func(tx *sql.Tx) error {
		logrus.Info("delete user")
		if err := s.UserService.Delete(ctx, tx, ID); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	logrus.Info("user deleted successfully")

	return nil
}

func (s *UserUseCase) Logout(refreshBearerToken string) error {
	ctx := context.Background()

	if err := middleware.WithTransaction(ctx, s.DB, func(tx *sql.Tx) error {
		logrus.Info("validating refresh token")

		refreshToken := strings.Split(refreshBearerToken, " ")[1]

		token, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) {
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
			return fmt.Errorf("access token expired")
		}

		logrus.Info("getting refresh token from database")

		refreshDatabaseToken, err := s.TokenService.GetTokenByRefresh(refreshToken, tx, ctx)
		if err != nil {
			return err
		}

		logrus.Info("delete refresh token from database")

		if err := s.TokenService.DeleteToken(refreshDatabaseToken.UserID, tx, ctx); err != nil {
			return err
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
