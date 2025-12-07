package service

import (
	"errors"

	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"github.com/nitikhon/golang-inventory-system/internal/core/port"
	"github.com/nitikhon/golang-inventory-system/internal/util"
	"gorm.io/gorm"
)

// Ensure UserService implements UserServiceInterface
var _ UserServiceInterface = (*UserService)(nil)

type UserServiceInterface interface {
	CreateUser(user *entity.User) (*entity.User, error)
	UpdateUser(user *entity.User) (*entity.User, error)
	DeleteUser(userID uint) error
	GetAllUsers() ([]*entity.User, error)
	GetUserByID(userID uint) (*entity.User, error)
	GetUserByUsername(username string) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByPhone(phone string) (*entity.User, error)
	Login(username, password string) (string, string, error)
	RefreshToken(refreshToken string) (entity.Token, error)
	Logout(userID uint) error
	UpdateUserProfile(user *entity.User) (*entity.User, error)
	UpdateUserPassword(user *entity.User) error
	UpdateUserEmail(user *entity.User) error
	UpdateUserAdminStatus(user *entity.User) error
}

// UserService provides the use cases for the user entity.
type UserService struct {
	repo   port.UserRepository
	crypto util.CryptoUtil
	jwt    util.JWTUtil
}

// NewUserService creates a new UserService instance.
func NewUserService(repo port.UserRepository, crypto util.CryptoUtil, jwt util.JWTUtil) *UserService {
	return &UserService{repo: repo, crypto: crypto, jwt: jwt}
}

// Create creates a new user.
func (s *UserService) CreateUser(user *entity.User) (*entity.User, error) {
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
	accessToken, err := s.jwt.GenerateAccessToken(*user)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.jwt.GenerateRefreshToken(*user)
	if err != nil {
		return "", "", err
	}

	// Save refresh token
	err = s.repo.UpdateRefreshToken(user.ID, refreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// RefreshToken validates the provided refresh token, retrieves the associated user,
// and generates a new pair of access and refresh tokens.
func (s *UserService) RefreshToken(refreshToken string) (entity.Token, error) {
	// Validate the Refresh Token
	userID, err := s.jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return entity.Token{}, err
	}

	// Retrieve the user
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return entity.Token{}, err
	}

	if user.RefreshToken != refreshToken {
		return entity.Token{}, errors.New("invalid refresh token")
	}

	// Generate new tokens
	accessToken, err := s.jwt.GenerateAccessToken(*user)
	if err != nil {
		return entity.Token{}, err
	}

	newRefreshToken, err := s.jwt.GenerateRefreshToken(*user)
	if err != nil {
		return entity.Token{}, err
	}

	// Save new refresh token
	err = s.repo.UpdateRefreshToken(user.ID, newRefreshToken)
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

	// Update the user in the repository
	err = s.repo.UpdateRefreshToken(user.ID, "")
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) UpdateUserProfile(user *entity.User) (*entity.User, error) {
	updatedFields := map[string]any{
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"phone":      user.Phone,
	}

	_, err := s.repo.UpdateUserProfile(user.ID, updatedFields)
	if err != nil {
		return &entity.User{}, err
	}

	updatedUser, err := s.repo.GetUserByID(user.ID)
	if err != nil {
		return &entity.User{}, errors.New("error while trying to get updated user")
	}

	return updatedUser, nil
}

func (s *UserService) UpdateUserPassword(user *entity.User) error {
	hashedPassword, err := s.crypto.HashPassword(user.Password)
	if err != nil {
		return err
	}

	err = s.repo.UpdateUserPassword(user.ID, hashedPassword)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) UpdateUserEmail(user *entity.User) error {
	if err := s.repo.UpdateUserEmail(user.ID, user.Email); err != nil {
		return err
	}
	return nil
}

func (s *UserService) UpdateUserAdminStatus(user *entity.User) error {
	if err := s.repo.UpdateUserAdminStatus(user.ID, user.IsAdmin); err != nil {
		return err
	}
	return nil
}
