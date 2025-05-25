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

func ComparePassword(hashedPassword string, plainPassword []byte) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), plainPassword)
	if err != nil {
		return err
	}
	return nil
}
