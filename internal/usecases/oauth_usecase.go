package usecases

import (
	"auth/internal/config"
	"auth/internal/domain/entity"
	"auth/internal/repository"
	"auth/pkg/helper"
	"auth/pkg/middleware"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

func (s *UserUseCase) GoogleAuth(code string) (*entity.JWTResponse, error) {
	ctx := context.Background()

	var jwtResp *entity.JWTResponse
	if err := middleware.WithTransaction(ctx, s.DB, func(tx *sql.Tx) error {
		token, err := config.OauthGoogleConfig.Exchange(ctx, code)
		if err != nil {
			return err
		}

		client := config.OauthGoogleConfig.Client(ctx, token)
		resp, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
		if err != nil {
			return err
		}

		var user *entity.OAuthUserData
		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			return err
		}
		defer resp.Body.Close()

		userResp, err := s.UserService.GetByEmail(ctx, tx, user.Email)
		if err == nil {
			userResp.Provider = "Google"
			userResp.ProviderID = user.ProviderID
			if err := s.UserService.Update(ctx, tx, userResp); err != nil {
				return err
			}
			logrus.Info("updating user data from google oauth")

			accessToken, err := helper.GenerateToken(userResp.ID, "user", 1, 0, 0)
			if err != nil {
				return err
			}

			refreshToken, err := helper.GenerateToken(userResp.ID, "user", 0, 3, 0)
			if err != nil {
				return err
			}

			token, err := s.TokenService.GetTokenByUserID(userResp.ID, tx, ctx)
			if err != nil {
				return err
			} else if token == nil {
				err = s.TokenService.SaveToken(&entity.RefreshToken{
					UserID:                userResp.ID,
					RefreshToken:          refreshToken,
					RefreshTokenExpiredAt: time.Now().AddDate(0, 3, 0),
				}, tx, ctx)
				if err != nil {
					return err
				}
				logrus.Info("refresh token created successfully and saved to database")
			} else {
				if err := s.TokenService.UpdateToken(&entity.RefreshToken{
					ID:                    userResp.ID,
					UserID:                userResp.ID,
					RefreshToken:          refreshToken,
					RefreshTokenExpiredAt: time.Now().AddDate(0, 3, 0),
				}, tx, ctx); err != nil {
					return err
				}
				logrus.Info("refresh token updated successfully and saved to database")
			}

			jwtResp = &entity.JWTResponse{
				Token:                 fmt.Sprintf("Bearer %s", accessToken),
				TokenExpiredAt:        time.Now().AddDate(1, 0, 0).Format("2006-01-02 15:04:05"),
				RefreshToken:          fmt.Sprintf("Bearer %s", refreshToken),
				RefreshTokenExpiredAt: time.Now().AddDate(0, 3, 0).Format("2006-01-02 15:04:05"),
			}

			return nil
		} else if err != sql.ErrNoRows && err != repository.ErrUserNotFound {
			return err
		}
		return err

	}); err != nil {
		return nil, err
	}

	return jwtResp, nil
}

func (s *UserUseCase) JWTCreateAndResponseUserToken(userID int) error {
	accessToken, err := helper.GenerateToken(userID, "user", 1, 0, 0)
	if err != nil {
		return err
	}

	refreshToken, err := helper.GenerateToken(userID, "user", 0, 3, 0)
	if err != nil {
		return err
	}

	token, err := s.TokenService.GetTokenByUserID(userID, tx, ctx)
	if err != nil {
		return err
	} else if token == nil {
		err = s.TokenService.SaveToken(&entity.RefreshToken{
			UserID:                userID,
			RefreshToken:          refreshToken,
			RefreshTokenExpiredAt: time.Now().AddDate(0, 3, 0),
		}, tx, ctx)
		if err != nil {
			return err
		}
		logrus.Info("refresh token created successfully and saved to database")
	} else {
		if err := s.TokenService.UpdateToken(&entity.RefreshToken{
			ID:                    userID,
			UserID:                userID,
			RefreshToken:          refreshToken,
			RefreshTokenExpiredAt: time.Now().AddDate(0, 3, 0),
		}, tx, ctx); err != nil {
			return err
		}
		logrus.Info("refresh token updated successfully and saved to database")
	}

	jwtResp = &entity.JWTResponse{
		Token:                 fmt.Sprintf("Bearer %s", accessToken),
		TokenExpiredAt:        time.Now().AddDate(1, 0, 0).Format("2006-01-02 15:04:05"),
		RefreshToken:          fmt.Sprintf("Bearer %s", refreshToken),
		RefreshTokenExpiredAt: time.Now().AddDate(0, 3, 0).Format("2006-01-02 15:04:05"),
	}
}
