package controller

import (
	"auth/internal/domain/entity"
	"auth/internal/domain/interface"
	"auth/internal/services"
	"auth/pkg/helper"
	"net/http"

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

	user := router.Group("/user")
	user.Use(authMiddleware)
	{
		user.POST("/logout", c.Logout)
		user.GET("/profile", c.GetProfile)
		user.DELETE("/delete", c.DeleteAccount)
	}
}

// Register handles user registration
func (c *UserController) Register(ctx *gin.Context) {
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

// GetProfile handles getting user profile
func (c *UserController) GetProfile(ctx *gin.Context) {
	logrus.Info("handling get profile request")

	id, exists := ctx.Get("userID")
	logrus.Error("handling get profile request", id)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	userID, err := helper.ChangeID(id)
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

	ctx.JSON(http.StatusOK, gin.H{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
	})
}

// DeleteAccount handles user account deletion
func (c *UserController) DeleteAccount(ctx *gin.Context) {
	// Get user ID from context (set by auth middleware)
	paramID, exists := ctx.Get("userID")

	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	logrus.Info("change type of user ID")
	userID, err := helper.ChangeID(paramID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is not valid"})
		return
	}

	logrus.Info("delete user")
	err = c.userUsecase.DeleteUser(userID)
	if err != nil {
		logrus.Error("error deleting user account: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "account deleted successfully"})
}

// Logout handles user logout
func (c *UserController) Logout(ctx *gin.Context) {
	logrus.Info("handling logout request")
	paramID, exists := ctx.Get("userID")
	logrus.Error("handling logout request", paramID)
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	logrus.Info("change type of user ID")
	userID, err := helper.ChangeID(paramID)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "user ID is not valid"})
		return
	}

	if err := c.userUsecase.Logout(userID); err != nil {
		logrus.Error("error logging out: ", err)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "logout successful"})
}
