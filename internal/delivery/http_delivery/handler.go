package httpdelivery

import (
	"auth/internal/config"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func NewHandler() *gin.Engine {
	router := gin.Default()

	logrus.Info("setup handler")

	config.Connect()

	return router
}
