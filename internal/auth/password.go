package auth

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func GeneratePasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		if errors.Is(err, bcrypt.ErrPasswordTooLong) {
			return "", ErrPasswordTooLong
		}
		return "", fmt.Errorf("generate hash from password: %w", err)
	}
	return string(hash), nil
}

func CheckPassword(hash string, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrWrongLoginOrPassword
		}
		return fmt.Errorf("compare hash and password: %w", err)
	}

	return nil
}
