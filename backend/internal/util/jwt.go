package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"github.com/nitikhon/golang-inventory-system/internal/util/errormap"
)

type JWTUtil interface {
	GenerateAccessToken(user entity.User) (string, error)
	GenerateRefreshToken(user entity.User) (string, error)
	ValidateAccessToken(tokenStr string) (uint, error)
	ValidateRefreshToken(tokenStr string) (uint, error)
}

type JWT struct{}

func NewAppJWTUtil() JWTUtil {
	return JWT{}
}

// GenerateAccessToken generates an access token for the given user.
// The access token is valid for 15 minutes.
func (s JWT) GenerateAccessToken(user entity.User) (string, error) {
	// Create the claims for the access token
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"is_admin": user.IsAdmin,
		"exp":      time.Now().Add(15 * time.Minute).Unix(),
	}

	// Create a new token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Set the audience and issuer claims
	signedToken, err := token.SignedString(GetEnvOrPanic("ACCESS_TOKEN_SECRET"))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// GenerateRefreshToken generates a refresh token for the given user.
// The refresh token is valid for 7 days.
func (s JWT) GenerateRefreshToken(user entity.User) (string, error) {
	// Create the claims for the refresh token
	claims := jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	// Create a new token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Set the audience and issuer claims
	signedToken, err := token.SignedString(GetEnvOrPanic("REFRESH_TOKEN_SECRET"))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ValidateAccessToken validates the access token and returns the user's data if valid.
func (s JWT) ValidateAccessToken(tokenStr string) (uint, error) {
	return validateToken(tokenStr, GetEnvOrPanic("ACCESS_TOKEN_SECRET"))
}

// ValidateRefreshToken validates the refresh token and returns the user ID if valid.
func (s JWT) ValidateRefreshToken(tokenStr string) (uint, error) {
	return validateToken(tokenStr, GetEnvOrPanic("REFRESH_TOKEN_SECRET"))
}

// validateToken validates the token and returns the user ID if valid.
func validateToken(tokenStr string, secret []byte) (uint, error) {
	// Parse the token
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})

	// Check if the token is valid
	if err != nil || !token.Valid {
		return 0, errors.New(errormap.ErrInvalidToken)
	}

	// Check the claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || claims["user_id"] == nil {
		return 0, errors.New(errormap.ErrInvalidClaims)
	}

	// Get the user ID from the claims
	// The user ID is stored as a float64 in the claims, so we need to convert it to int64
	userID := uint(claims["user_id"].(float64))
	return userID, nil
}
