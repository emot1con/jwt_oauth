package httpdelivery

import (
	"auth/internal/config"
	"auth/internal/controller"
	"auth/internal/repository"
	"auth/internal/services"
	"auth/internal/usecases"
	"auth/pkg/middleware"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

func NewHandler() *gin.Engine {
	router := gin.Default()

	validator := validator.New()

	DB := config.Connect()

	if DB == nil {
		logrus.Panic("failed to connect to database")
	}

	config.InitRedis()

	userRepo := repository.NewUserRepository()
	tokenRepo := repository.NewTokenRepository()

	userService := services.NewUserService(userRepo)
	tokenService := services.NewTokenService(tokenRepo)

	userUsecase := usecases.NewUserUseCase(userService, tokenService, DB, validator)

	userController := controller.NewUserController(userUsecase, tokenService)
	userController.RegisterRoutes(router, middleware.ProtectedEndpoint())

	return router
}
