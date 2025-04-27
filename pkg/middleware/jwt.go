package middleware

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthWithJWT(bearerToken string) gin.HandlerFunc {
	return func(c *gin.Context) {
		stringToken := strings.Split(bearerToken, " ")[1]
		token, err := jwt.ParseWithClaims(stringToken, &jwt.RegisteredClaims{}, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(os.Getenv("SECRET_KEY")), nil
		})
		if err != nil || !token.Valid {
			c.JSON(401, gin.H{"error": "invalid token"})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(401, gin.H{"error": "invalid claims"})
		}

		expAt := int64(claims["exp"].(float64))
		if time.Now().Unix() > expAt {
			c.JSON(401, gin.H{"error": "token expired"})
		}

		ID := int(claims["id"].(float64))
		c.Set("userID", ID)
		role := claims["role"].(string)
		c.Set("role", role)
	}
}
