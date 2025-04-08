package service

import (
	"errors"

	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"github.com/nitikhon/golang-inventory-system/internal/core/port"
	"github.com/nitikhon/golang-inventory-system/internal/util"
)

type AuthService struct {
	repo port.UserRepository
}

func NewAuthService(repo port.UserRepository) *AuthService {
	return &AuthService{repo: repo}
}

// Register creates a new user with the given email and password.
func (s *AuthService) Register(user entity.User) (*entity.User, error) {
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return &entity.User{}, err
	}

	user.Password = hashedPassword

	return s.repo.CreateUser(&user)
}

// Login checks if the user exists and if the password is correct.
func (s *AuthService) Login(username, password string) (string, string, error) {
	// Check if user exists
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return "", "", err
	}

	// Check if password is correct
	if err := util.CheckPasswordHash(user.Password, password); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	// Generate access and refresh tokens
	access, err := util.GenerateAccessToken(*user)
	if err != nil {
		return "", "", err
	}

	refresh, err := util.GenerateRefreshToken(*user)
	if err != nil {
		return "", "", err
	}

	// Save refresh token
	user.RefreshToken = refresh
	if err := s.repo.UpdateRefreshToken(user.ID, refresh); err != nil {
		return "", "", err
	}

	return access, refresh, nil
}

func (s *AuthService) RefreshToken(refreshToken string) (struct {
    AccessToken  string
    RefreshToken string
}, error) {
    // Validate the Refresh Token
    userID, err := util.ValidateRefreshToken(refreshToken)
    if err != nil {
        return struct {
            AccessToken  string
            RefreshToken string
        }{}, err
    }

    // Retrieve the user
    user, err := s.repo.GetUserByID(userID)
    if err != nil {
        return struct {
            AccessToken  string
            RefreshToken string
        }{}, err
    }

    // Generate new tokens
    accessToken, err := util.GenerateAccessToken(*user)
    if err != nil {
        return struct {
            AccessToken  string
            RefreshToken string
        }{}, err
    }

    refreshToken, err = util.GenerateRefreshToken(*user)
    if err != nil {
        return struct {
            AccessToken  string
            RefreshToken string
        }{}, err
    }

    return struct {
        AccessToken  string
        RefreshToken string
    }{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
    }, nil
}

// GetUserByID retrieves a user by their ID.
func (s *AuthService) GetUserByID(id uint) (*entity.User, error) {
	user, err := s.repo.GetUserByID(id)
	if err != nil {
		return &entity.User{}, err
	}

	return user, nil
}
