package util

import (
	"golang.org/x/crypto/bcrypt"
)

type CryptoUtil interface {
	HashPassword(password string) (string, error)
	CheckPasswordHash(hashedPassword, password string) error
}

type AppCryptoUtil struct {}

func NewAppCrptoUtil() AppCryptoUtil {
	return AppCryptoUtil{}
}

// HashPassword hashes the password using bcrypt.
func (u AppCryptoUtil) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

// CheckPasswordHash checks if the provided password matches the hashed password.
func (u AppCryptoUtil) CheckPasswordHash(hashedPassword, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return err
	}
	return nil
}