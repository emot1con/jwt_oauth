package controller

import (
	"auth/internal/domain/entity"
	interfaces "auth/internal/domain/interface"
	"auth/internal/services"
	"auth/pkg/helper"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type UserController struct {
	userUsecase  interfaces.UserUsecaseInterface
	tokenService *services.TokenService
}

func NewUserController(userUsecase interfaces.UserUsecaseInterface, tokenService *services.TokenService) *UserController {
	return &UserController{
		userUsecase:  userUsecase,
		tokenService: tokenService,
	}
}

// RegisterRoutes registers all auth routes
func (c *UserController) RegisterRoutes(router *gin.Engine, authMiddleware gin.HandlerFunc) {
	auth := router.Group("/auth")
	{
		auth.POST("/register", c.Register)
		auth.POST("/login", c.Login)
	}
	oauth := router.Group("/oauth")
	{
		oauth.GET("/google", c.OAuthGoogle)
		oauth.GET("/facebook", c.OAuthFacebook)
		oauth.GET("/github", c.OAuthGithub)

		oauth.GET("/google/callback", c.OAuthGoogleCallback)
		oauth.GET("/facebook/callback", c.OAuthFacebookCallback)
		oauth.GET("/github/callback", c.OAuthGithubCallback)
	}

	user := router.Group("/user")
	user.Use(authMiddleware)
	{
		user.POST("/logout", c.Logout)
		user.GET("/profile", c.GetProfile)
		user.DELETE("/delete", c.DeleteAccount)
		user.POST("/refresh", c.RefreshToken)
	}
}

// Register handles user registration
func (c *UserController) Register(ctx *gin.Context) {
	isLimit, err := helper.RateLimiter(ctx, "register_req")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !isLimit {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
		return
	}

	logrus.Info("handling register request")

	var payload entity.RegisterPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.userUsecase.Register(&payload)
	if err != nil {
		logrus.Error("error registering user: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

// Login handles user login
func (c *UserController) Login(ctx *gin.Context) {

	isLimit, err := helper.RateLimiter(ctx, "login_req")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !isLimit {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
		return
	}

	logrus.Info("handling login request")

	var payload entity.LoginPayload
	if err := ctx.ShouldBindJSON(&payload); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := c.userUsecase.Login(&payload)
	if err != nil {
		logrus.Error("error logging in: ", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *UserController) RefreshToken(ctx *gin.Context) {

	isLimit, err := helper.RateLimiter(ctx, "refresh_token_req")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !isLimit {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
		return
	}

	logrus.Info("handling refresh token request")

	id, exist := ctx.Get("userID")
	if !exist {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := helper.FormatIDToInt(id)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not valid user ID"})
		return

	}

	header := ctx.GetHeader("Authorization")
	refreshToken := strings.Split(header, " ")
	refreshtokenResponse, err := c.userUsecase.RefreshToken(refreshToken[1], userID)

	if err != nil {
		logrus.Error("error refreshing token: ", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, refreshtokenResponse)
}

// GetProfile handles getting user profile
func (c *UserController) GetProfile(ctx *gin.Context) {

	isLimit, err := helper.RateLimiter(ctx, "get_profile_req")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !isLimit {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
		return
	}

	logrus.Info("handling get profile request")

	id, exists := ctx.Get("userID")
	logrus.Error("handling get profile request", id)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := helper.FormatIDToInt(id)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "not valid user ID"})
		return
	}

	user, err := c.userUsecase.GetUserByID(userID)
	if err != nil {
		logrus.Error("error getting user profile: ", err)
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	logrus.Info("user profile founded")

	ctx.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

// DeleteAccount handles user account deletion
func (c *UserController) DeleteAccount(ctx *gin.Context) {

	isLimit, err := helper.RateLimiter(ctx, "delete_account_req")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !isLimit {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
		return
	}

	// Get user ID from context (set by auth middleware)
	header := ctx.GetHeader("Authorization")
	token := strings.Split(header, " ")

	paramID, exists := ctx.Get("userID")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	logrus.Info("change type of user ID")
	userID, err := helper.FormatIDToInt(paramID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is not valid"})
		return
	}

	logrus.Info("delete user")
	err = c.userUsecase.DeleteUser(userID, token[1])
	if err != nil {
		logrus.Error("error deleting user account: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "account deleted successfully"})
}

// Logout handles user logout
func (c *UserController) Logout(ctx *gin.Context) {
	isLimit, err := helper.RateLimiter(ctx, "logout_req")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if !isLimit {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "too many requests"})
		return
	}

	logrus.Info("handling logout request")
	header := ctx.GetHeader("Authorization")
	token := strings.Split(header, " ")

	paramID, exists := ctx.Get("userID")
	logrus.Error("handling logout request", paramID)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	logrus.Info("change type of user ID")
	userID, err := helper.FormatIDToInt(paramID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is not valid"})
		return
	}

	if err := c.userUsecase.Logout(userID, token[1]); err != nil {
		logrus.Error("error logging out: ", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
