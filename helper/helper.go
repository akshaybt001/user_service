package helper

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(pass string) (string, error) {

	password := []byte(pass)

	hash, err := bcrypt.GenerateFromPassword(password, 10)
	if err != nil {
		return "", fmt.Errorf("can't hash password")
	}
	return string(hash), nil
}

func VerifyPassword(password, checkpassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(password), []byte(checkpassword))
}
