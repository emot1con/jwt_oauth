package controller

import (
	"auth/internal/config"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

func (c *UserController) OAuthGoogle(ctx *gin.Context) {
	url := config.OauthGoogleConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	ctx.Redirect(http.StatusFound, url)
	logrus.Info("redirecting to google oauth")
}

func (c *UserController) OAuthFacebook(ctx *gin.Context) {
	url := config.OauthFacebookConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	ctx.Redirect(http.StatusFound, url)
	logrus.Info("redirecting to facebook oauth")
}

func (c *UserController) OAuthGithub(ctx *gin.Context) {
	url := config.OauthGithubConfig.AuthCodeURL("state", oauth2.AccessTypeOffline)
	ctx.Redirect(http.StatusFound, url)
	logrus.Info("redirecting to github oauth")
}

func (c *UserController) OAuthGoogleCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	if code == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "missing code"})
		return
	}

	logrus.Info("usecase google auth")
	jwtToken, err := c.userUsecase.GoogleAuth(code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	logrus.Infof("getting client info from google successful with jwt: %s", jwtToken)

	ctx.JSON(http.StatusOK, jwtToken)
	logrus.Info("google oauth callback successful")
}
