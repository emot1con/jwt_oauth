package middleware

import (
	"auth/internal/services"
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

/// AuthWithJWT is a middleware function that checks the JWT token in the request header.
// This middlware ware checks if the token is valid and not expired. If valid, it sets the user ID and role in the context.

type contextKey string

var UserKey contextKey = "userID"

func ProtectedEndpoint(s services.TokenService) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")

		parts := strings.Split(header, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token missing"})
			c.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

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

		userIDClaims := claims["userID"].(string)
		userID, err := strconv.Atoi(userIDClaims)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Can't convert userID to int"})
			c.Abort()
			return
		}

		_, err = s.GetTokenByRefresh(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}

		ctx := c.Request.Context()
		ctx = context.WithValue(ctx, UserKey, userID)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
