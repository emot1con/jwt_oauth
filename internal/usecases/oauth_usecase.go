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
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
)

func (s *UserUseCase) GoogleAuth(code string) (*entity.JWTResponse, error) {
	logrus.Info("handling Google OAuth authentication callback")

	ctx := context.Background()

	var jwtResp *entity.JWTResponse
	if err := middleware.WithTransaction(ctx, s.DB, func(tx *sql.Tx) error {
		token, err := config.OauthGoogleConfig.Exchange(ctx, code)
		if err != nil {
			logrus.Error("1 error: ", err)
			return err
		}

		client := config.OauthGoogleConfig.Client(ctx, token)
		resp, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo")
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		oauthUser := &entity.OAuthUserData{}
		if err := json.NewDecoder(resp.Body).Decode(oauthUser); err != nil {
			return err
		}

		if oauthUser.Email == "" {
			logrus.Error("email not found in google oauth response")
			return fmt.Errorf("email not found in google oauth response")
		}

		oauthUser.Provider = "Google"

		userResp, err := s.UserService.GetByEmail(ctx, tx, oauthUser.Email)
		if err == nil {
			userResp.Provider = "Google"
			userResp.ProviderID = oauthUser.ProviderID

			if err := s.UserService.Update(ctx, tx, userResp); err != nil {
				return err
			}
			logrus.Info("updating user data from google oauth: ", userResp.Email, userResp.Provider)

			jwtResp, err = s.JWTCreateAndResponseUserOauthToken(userResp.ID, tx, ctx)
			if err != nil {
				return err
			}

			return nil
		} else if err != sql.ErrNoRows && err != repository.ErrUserNotFound {
			return err
		}

		logrus.Info("creating new user from google oauth", userResp.Email, userResp.Provider)
		jwtResp, err = s.CreateUserWithResponse(oauthUser, tx, ctx)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return jwtResp, nil
}

func (s *UserUseCase) GitHubAuth(code string) (*entity.JWTResponse, error) {
	logrus.Info("handling Github OAuth authentication callback")

	ctx := context.Background()

	var jwtResp *entity.JWTResponse
	if err := middleware.WithTransaction(ctx, s.DB, func(tx *sql.Tx) error {
		token, err := config.OauthGithubConfig.Exchange(ctx, code)
		if err != nil {
			logrus.Error("1 error: ", err)
			return err
		}

		client := config.OauthGithubConfig.Client(ctx, token)
		resp, err := client.Get("https://api.github.com/user")
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		oauthUser := &entity.OauthGithubUserModel{}
		if err := json.NewDecoder(resp.Body).Decode(oauthUser); err != nil {
			return err
		}

		if oauthUser.Email == "" {
			logrus.Error("email not found in github oauth response")
			return fmt.Errorf("email not found")
		}

		logrus.Info("update user from github oauth", oauthUser)

		oauthUser.Provider = "GitHub"
		userResp, err := s.UserService.GetByEmail(ctx, tx, oauthUser.Email)
		if err == nil {
			logrus.Info("update user from github oauth", userResp)

			userResp.Provider = "GitHub"
			userResp.ProviderID = strconv.Itoa(int(oauthUser.ProviderID))

			if err := s.UserService.Update(ctx, tx, userResp); err != nil {
				return err
			}

			jwtResp, err = s.JWTCreateAndResponseUserOauthToken(userResp.ID, tx, ctx)
			if err != nil {
				return err
			}

			return nil
		} else if err != sql.ErrNoRows && err != repository.ErrUserNotFound {
			return err
		}

		parsedOauthUserModel := &entity.OAuthUserData{
			ProviderID: strconv.Itoa(int(oauthUser.ProviderID)),
			Provider:   oauthUser.Provider,
			Email:      oauthUser.Email,
			Name:       oauthUser.Login,
			AvatarURL:  oauthUser.AvatarURL,
		}

		logrus.Info("creating new user from Github oauth", parsedOauthUserModel.Email, parsedOauthUserModel.Provider)
		jwtResp, err = s.CreateUserWithResponse(parsedOauthUserModel, tx, ctx)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return jwtResp, nil
}

func (s *UserUseCase) FacebookAuth(code string) (*entity.JWTResponse, error) {
	logrus.Info("handling Facebook OAuth authentication callback")

	ctx := context.Background()

	var jwtResp *entity.JWTResponse
	if err := middleware.WithTransaction(ctx, s.DB, func(tx *sql.Tx) error {
		token, err := config.OauthFacebookConfig.Exchange(ctx, code)
		if err != nil {
			logrus.Error("1 error: ", err)
			return err
		}

		client := config.OauthFacebookConfig.Client(ctx, token)
		resp, err := client.Get("https://graph.facebook.com/me?fields=id,name,email&access_token=" + token.AccessToken)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		oauthUser := &entity.OAuthUserData{}
		if err := json.NewDecoder(resp.Body).Decode(oauthUser); err != nil {
			return err
		}

		if oauthUser.Email == "" {
			return fmt.Errorf("email not found")
		}

		oauthUser.Provider = "Facebook"
		logrus.Info("update user from Facebook oauth1", oauthUser)

		userResp, err := s.UserService.GetByEmail(ctx, tx, oauthUser.Email)
		if err != sql.ErrNoRows && err != repository.ErrUserNotFound {
			logrus.Error("error get user by email 1: ", err)
			return err
		} else if userResp != nil {
			logrus.Info("update user from Facebook oauth2", userResp)

			userResp.Provider = "Facebook"
			userResp.ProviderID = oauthUser.ProviderID

			if err := s.UserService.Update(ctx, tx, userResp); err != nil {
				return err
			}

			jwtResp, err = s.JWTCreateAndResponseUserOauthToken(userResp.ID, tx, ctx)
			if err != nil {
				return err
			}

			return nil
		}

		jwtResp, err = s.CreateUserWithResponse(oauthUser, tx, ctx)
		if err != nil {
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return jwtResp, nil
}

func (s *UserUseCase) CreateUserWithResponse(oauthUser *entity.OAuthUserData, tx *sql.Tx, ctx context.Context) (*entity.JWTResponse, error) {
	logrus.Info("creating new user from oauth", oauthUser.Email, oauthUser.Provider, oauthUser)
	if err := s.UserService.Create(ctx, tx, &entity.RegisterPayload{
		Name:     oauthUser.Name,
		Email:    oauthUser.Email,
		Password: "",
		Provider: oauthUser.Provider,
	}); err != nil {
		return nil, err
	}
	logrus.Info("succedd creating new user from oauth")

	createdUser, err := s.UserService.GetByEmail(ctx, tx, oauthUser.Email)
	if err != nil {
		return nil, err
	}

	jwtResp, err := s.JWTCreateAndResponseUserOauthToken(createdUser.ID, tx, ctx)
	if err != nil {
		return nil, err
	}

	return jwtResp, nil
}

func (s *UserUseCase) JWTCreateAndResponseUserOauthToken(userID int, tx *sql.Tx, ctx context.Context) (*entity.JWTResponse, error) {
	accessToken, err := helper.GenerateToken(userID, "user", 1, 0, 0)
	if err != nil {
		return nil, err
	}

	refreshToken, err := helper.GenerateToken(userID, "user", 0, 3, 0)
	if err != nil {
		return nil, err
	}

	token, err := s.TokenService.GetTokenByUserID(userID, tx, ctx)
	if err != nil {
		return nil, err
	} else if token == nil {
		err = s.TokenService.SaveToken(&entity.RefreshToken{
			UserID:                userID,
			RefreshToken:          refreshToken,
			RefreshTokenExpiredAt: time.Now().AddDate(0, 3, 0),
		}, tx, ctx)
		if err != nil {
			return nil, err
		}
		logrus.Info("refresh token created successfully and saved to database")
	} else {
		if err := s.TokenService.UpdateToken(&entity.RefreshToken{
			ID:                    userID,
			UserID:                userID,
			RefreshToken:          refreshToken,
			RefreshTokenExpiredAt: time.Now().AddDate(0, 3, 0),
		}, tx, ctx); err != nil {
			return nil, err
		}
		logrus.Info("refresh token updated successfully and saved to database")
	}

	return &entity.JWTResponse{
		Token:                 fmt.Sprintf("Bearer %s", accessToken),
		TokenExpiredAt:        time.Now().AddDate(1, 0, 0).Format("2006-01-02 15:04:05"),
		RefreshToken:          fmt.Sprintf("Bearer %s", refreshToken),
		RefreshTokenExpiredAt: time.Now().AddDate(0, 3, 0).Format("2006-01-02 15:04:05"),
	}, nil

}
