package middleware

import (
	"auth/pkg/helper"
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type contextKey string

var UserKey contextKey = "userID"

// / ProtectedEnpoint is a middleware function that checks the JWT token in the request header.
// This middlware ware checks if the token is valid and not expired. If valid, it sets the user ID and role in the context.
func ProtectedEndpoint() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")

		logrus.Info("checking authorization header")

		parts := strings.Split(header, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token missing"})
			c.Abort()
			return
		}

		logrus.Info("Cut authorization header")

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(os.Getenv("SECRET_KEY")), nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		blacklisted, _ := helper.IsTokenBlacklisted(tokenString)
		if blacklisted {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token has been revoked"})
			c.Abort()
			return
		}

		logrus.Info("parsing token")

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid claims"})
			c.Abort()
			return
		}

		expiredAt := int64(claims["exp"].(float64))
		if time.Now().Unix() > expiredAt {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}

		userID := int64(claims["user_id"].(float64))

		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, UserKey, userID)
		c.Request = c.Request.WithContext(ctx)
		c.Set("userID", userID)

		c.Next()
	}
}
