package helper

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateToken(id int, role string, days, months, years int) (string, error) {
	expirationTime := time.Now().AddDate(years, months, days).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"exp":     expirationTime,
		"role":    role,
	})

	return token.SignedString([]byte(os.Getenv("SECRET_KEY")))
}
