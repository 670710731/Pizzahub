package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidPassword = errors.New("password must be 8-72 bytes")

func HashPassword(password string) (string, error) {
	if len(password) < 8 || len(password) > 72 {
		return "", ErrInvalidPassword
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func VerifyPassword(hash, password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
