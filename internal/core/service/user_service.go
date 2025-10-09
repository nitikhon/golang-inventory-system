package service

import (
	"errors"

	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"github.com/nitikhon/golang-inventory-system/internal/core/port"
	"github.com/nitikhon/golang-inventory-system/internal/util"
	"github.com/nitikhon/golang-inventory-system/internal/util/validation"
	"gorm.io/gorm"
)

// Ensure UserService implements UserServiceInterface
var _ UserServiceInterface = (*UserService)(nil)

type UserServiceInterface interface {
	CreateUser(user *entity.User) (*entity.User, error)
	UpdateUser(user *entity.User) (*entity.User, error)
	DeleteUser(id uint) error
	GetAllUsers() ([]*entity.User, error)
	GetUserByID(id uint) (*entity.User, error)
	GetUserByUsername(username string) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByPhone(phone string) (*entity.User, error)
	UpdateRefreshToken(userID uint, refreshToken string) (*entity.User, error)
	Login(username, password string) (string, string, error)
	RefreshToken(refreshToken string) (entity.Token, error)
	Logout(userID uint) error
}

// UserService provides the use cases for the user entity.
type UserService struct {
	repo   port.UserRepository
	crypto util.CryptoUtil
}

// NewUserService creates a new UserService instance.
func NewUserService(repo port.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Create creates a new user.
func (s *UserService) CreateUser(user *entity.User) (*entity.User, error) {
	// Validate user fields
	_, err := validation.ValidateAndNormalizeUser(user)
	if err != nil {
		return &entity.User{}, err
	}

	// Check if username already exists
	existingUser, err := s.repo.GetUserByUsername(user.Username)
	if err == nil && existingUser != nil {
		return &entity.User{}, errors.New("a user with the provided credentials already exists")
	}

	// Check if email already exists
	existingUser, err = s.repo.GetUserByEmail(user.Email)
	if err == nil && existingUser != nil {
		return &entity.User{}, errors.New("a user with the provided credentials already exists")
	}

	// Check if phone already exists
	existingUser, err = s.repo.GetUserByPhone(user.Phone)
	if err == nil && existingUser != nil {
		return &entity.User{}, errors.New("a user with the provided credentials already exists")
	}

	// Hash the password before saving
	hashedPassword, err := s.crypto.HashPassword(user.Password)
	if err != nil {
		return &entity.User{}, err
	}

	user.Password = hashedPassword

	return s.repo.CreateUser(user)
}

// Update updates an existing user.
func (s *UserService) UpdateUser(user *entity.User) (*entity.User, error) {
	// Validate user fields
	_, err := validation.ValidateAndNormalizeUser(user)
	if err != nil {
		return &entity.User{}, err
	}

	return s.repo.UpdateUser(user)
}

// Delete deletes a user by their ID.
func (s *UserService) DeleteUser(id uint) error {
	return s.repo.DeleteUser(id)
}

// GetAllUsers returns all users.
func (s *UserService) GetAllUsers() ([]*entity.User, error) {
	return s.repo.GetAllUsers()
}

// GetUserByID returns a user by their ID.
func (s *UserService) GetUserByID(id uint) (*entity.User, error) {
	return s.repo.GetUserByID(id)
}

// GetUserByUsername retrieves a user by their username.
func (s *UserService) GetUserByUsername(username string) (*entity.User, error) {
	return s.repo.GetUserByUsername(username)
}

// GetUserByEmail retrieves a user by their email.
func (s *UserService) GetUserByEmail(email string) (*entity.User, error) {
	return s.repo.GetUserByEmail(email)
}

// GetUserByPhone retrieves a user by their phone number.
func (s *UserService) GetUserByPhone(phone string) (*entity.User, error) {
	return s.repo.GetUserByPhone(phone)
}

// UpdateRefreshToken updates the refresh token for a user.
func (s *UserService) UpdateRefreshToken(userID uint, refreshToken string) (*entity.User, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return &entity.User{}, err
	}

	if user == nil {
		return &entity.User{}, gorm.ErrRecordNotFound
	}

	if refreshToken == "" {
		return &entity.User{}, errors.New("refresh token cannot be empty")
	}

	user.RefreshToken = refreshToken

	return s.repo.UpdateUser(user)
}

// Login checks if the user exists and if the password is correct.
func (s *UserService) Login(username, password string) (string, string, error) {
	// Check if user exists
	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return "", "", err
	}

	if user == nil {
		return "", "", gorm.ErrRecordNotFound
	}

	// Check if password is correct
	if err := s.crypto.CheckPasswordHash(user.Password, password); err != nil {
		return "", "", errors.New("invalid credentials")
	}

	// Generate access and refresh tokens
	accessToken, err := util.GenerateAccessToken(*user)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := util.GenerateRefreshToken(*user)
	if err != nil {
		return "", "", err
	}

	user.RefreshToken = refreshToken

	// Save refresh token
	_, err = s.repo.UpdateUser(user)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// RefreshToken validates the provided refresh token, retrieves the associated user,
// and generates a new pair of access and refresh tokens.
func (s *UserService) RefreshToken(refreshToken string) (entity.Token, error) {
	// Validate the Refresh Token
	userID, err := util.ValidateRefreshToken(refreshToken)
	if err != nil {
		return entity.Token{}, err
	}

	// Retrieve the user
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return entity.Token{}, err
	}

	// Generate new tokens
	accessToken, err := util.GenerateAccessToken(*user)
	if err != nil {
		return entity.Token{}, err
	}

	newRefreshToken, err := util.GenerateRefreshToken(*user)
	if err != nil {
		return entity.Token{}, err
	}

	user.RefreshToken = newRefreshToken

	// Save new refresh token
	_, err = s.repo.UpdateUser(user)
	if err != nil {
		return entity.Token{}, err
	}

	return entity.Token{AccessToken: accessToken, RefreshToken: newRefreshToken}, nil
}

// Logout clears the refresh token for the user.
func (s *UserService) Logout(userID uint) error {
	// Retrieve the user
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return err
	}

	if user == nil {
		return gorm.ErrRecordNotFound
	}

	// Clear the refresh token
	user.RefreshToken = ""

	// Update the user in the repository
	_, err = s.repo.UpdateUser(user)
	if err != nil {
		return err
	}

	return nil
}
