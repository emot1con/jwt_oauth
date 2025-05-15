package delivery

import (
	"auth/internal/config"
	"auth/internal/controller"
	"auth/internal/repository"
	"auth/internal/services"
	"auth/internal/usecases"
	"auth/pkg/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

func NewHandler() *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 60 * 60,
	}))

	validator := validator.New()

	DB := config.Connect()

	if DB == nil {
		logrus.Panic("failed to connect to database")
	}

	config.InitRedis()
	logrus.Info("connected to redis")

	config.InitOauth()
	logrus.Info("oauth initialized")

	userRepo := repository.NewUserRepository()
	tokenRepo := repository.NewTokenRepository()

	userService := services.NewUserService(userRepo)
	tokenService := services.NewTokenService(tokenRepo)

	userUsecase := usecases.NewUserUseCase(userService, tokenService, DB, validator)

	userController := controller.NewUserController(userUsecase, tokenService)
	userController.RegisterRoutes(router, middleware.ProtectedEndpoint())

	return router
}
