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

	token, err := config.OauthGoogleConfig.Exchange(ctx, code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	client := config.OauthGoogleConfig.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v1/userinfo?alt=json")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	ctx.JSON(http.StatusOK, gin.H{"token": token})
	logrus.Info("google oauth callback successful")
}
