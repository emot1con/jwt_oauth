package helper

import (
	"golang.org/x/crypto/bcrypt"
)

func GenerateHashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}

// func ComparePassword(hashedPassword string, password []byte) bool {
// 	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), password); err != nil {
// 		logrus.Error("error comparing password: ", err)
// 		return false
// 	}

// 	return true
// }

func ComparePassword(hashedPassword string, plainPassword []byte) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), plainPassword)
	if err != nil {
		return err
	}
	return nil
}
