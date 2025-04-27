package helper

import "golang.org/x/crypto/bcrypt"

func GenerateHashPassword(password string) ([]byte, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}

func ComparePassword(password string, hashedPassword []byte) bool {
	if err := bcrypt.CompareHashAndPassword(hashedPassword, []byte(password)); err != nil {
		return false
	}

	return true
}
