package service

import (
	"context"
	"errors"

	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"github.com/nitikhon/golang-inventory-system/internal/core/port"
	"github.com/nitikhon/golang-inventory-system/internal/util"
	"github.com/nitikhon/golang-inventory-system/internal/util/errormap"
	"gorm.io/gorm"
)

// Ensure UserService implements UserServiceInterface
var _ UserServiceInterface = (*UserService)(nil)

type UserServiceInterface interface {
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	DeleteUser(ctx context.Context, userID uint) error
	GetAllUsers(ctx context.Context) ([]*entity.User, error)
	GetUserByID(ctx context.Context, userID uint) (*entity.User, error)
	GetUserByUsername(uctx context.Context, sername string) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*entity.User, error)
	Login(ctx context.Context, username, password string) (string, string, error)
	RefreshToken(ctx context.Context, refreshToken string) (entity.Token, error)
	Logout(ctx context.Context, userID uint) error
	UpdateUserProfile(ctx context.Context, user *entity.User) (*entity.User, error)
	UpdateUserPassword(ctx context.Context, user *entity.User) error
	UpdateUserEmail(ctx context.Context, user *entity.User) error
	UpdateUserAdminStatus(ctx context.Context, user *entity.User) error
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
func (s *UserService) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	// Check if username already exists
	existingUser, err := s.repo.GetUserByUsername(ctx, user.Username)
	if err == nil && existingUser != nil {
		return &entity.User{}, errors.New(errormap.ErrUserCredentialsExist)
	}

	// Check if email already exists
	existingUser, err = s.repo.GetUserByEmail(ctx, user.Email)
	if err == nil && existingUser != nil {
		return &entity.User{}, errors.New(errormap.ErrUserCredentialsExist)
	}

	// Check if phone already exists
	existingUser, err = s.repo.GetUserByPhone(ctx, user.Phone)
	if err == nil && existingUser != nil {
		return &entity.User{}, errors.New(errormap.ErrUserCredentialsExist)
	}

	// Hash the password before saving
	hashedPassword, err := s.crypto.HashPassword(user.Password)
	if err != nil {
		return &entity.User{}, err
	}

	user.Password = hashedPassword

	return s.repo.CreateUser(ctx, user)
}

// Update updates an existing user.
func (s *UserService) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	existingUser, err := s.repo.GetUserByEmail(ctx, user.Email)
	if err == nil && existingUser != nil && existingUser.ID != user.ID {
		return &entity.User{}, errors.New(errormap.ErrEmailAlreadyTaken)
	}

	existingUser, err = s.repo.GetUserByPhone(ctx, user.Phone)
	if err == nil && existingUser != nil && existingUser.ID != user.ID {
		return &entity.User{}, errors.New(errormap.ErrPhoneAlreadyTaken)
	}

	return s.repo.UpdateUser(ctx, user)
}

// Delete deletes a user by their ID.
func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	return s.repo.DeleteUser(ctx, id)
}

// GetAllUsers returns all users.
func (s *UserService) GetAllUsers(ctx context.Context) ([]*entity.User, error) {
	return s.repo.GetAllUsers(ctx)
}

// GetUserByID returns a user by their ID.
func (s *UserService) GetUserByID(ctx context.Context, id uint) (*entity.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

// GetUserByUsername retrieves a user by their username.
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	return s.repo.GetUserByUsername(ctx, username)
}

// GetUserByEmail retrieves a user by their email.
func (s *UserService) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	return s.repo.GetUserByEmail(ctx, email)
}

// GetUserByPhone retrieves a user by their phone number.
func (s *UserService) GetUserByPhone(ctx context.Context, phone string) (*entity.User, error) {
	return s.repo.GetUserByPhone(ctx, phone)
}

// Login checks if the user exists and if the password is correct.
func (s *UserService) Login(ctx context.Context, username, password string) (string, string, error) {
	// Check if user exists
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", "", err
	}

	if user == nil {
		return "", "", gorm.ErrRecordNotFound
	}

	// Check if password is correct
	if err := s.crypto.CheckPasswordHash(user.Password, password); err != nil {
		return "", "", errors.New(errormap.ErrInvalidCredentials)
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
	err = s.repo.UpdateRefreshToken(ctx, user.ID, refreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// RefreshToken validates the provided refresh token, retrieves the associated user,
// and generates a new pair of access and refresh tokens.
func (s *UserService) RefreshToken(ctx context.Context, refreshToken string) (entity.Token, error) {
	// Validate the Refresh Token
	userID, err := s.jwt.ValidateRefreshToken(refreshToken)
	if err != nil {
		return entity.Token{}, err
	}

	// Retrieve the user
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return entity.Token{}, err
	}

	if user.RefreshToken != refreshToken {
		return entity.Token{}, errors.New(errormap.ErrInvalidRefreshToken)
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
	err = s.repo.UpdateRefreshToken(ctx, user.ID, newRefreshToken)
	if err != nil {
		return entity.Token{}, err
	}

	return entity.Token{AccessToken: accessToken, RefreshToken: newRefreshToken}, nil
}

// Logout clears the refresh token for the user.
func (s *UserService) Logout(ctx context.Context, userID uint) error {
	// Retrieve the user
	user, err := s.repo.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}

	if user == nil {
		return gorm.ErrRecordNotFound
	}

	// Update the user in the repository
	err = s.repo.UpdateRefreshToken(ctx, user.ID, "")
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) UpdateUserProfile(ctx context.Context, user *entity.User) (*entity.User, error) {
	updatedFields := map[string]any{
		"first_name": user.FirstName,
		"last_name":  user.LastName,
		"phone":      user.Phone,
	}

	_, err := s.repo.UpdateUserProfile(ctx, user.ID, updatedFields)
	if err != nil {
		return &entity.User{}, err
	}

	updatedUser, err := s.repo.GetUserByID(ctx, user.ID)
	if err != nil {
		return &entity.User{}, errors.New(errormap.ErrGetUpdatedUser)
	}

	return updatedUser, nil
}

func (s *UserService) UpdateUserPassword(ctx context.Context, user *entity.User) error {
	hashedPassword, err := s.crypto.HashPassword(user.Password)
	if err != nil {
		return err
	}

	err = s.repo.UpdateUserPassword(ctx, user.ID, hashedPassword)
	if err != nil {
		return err
	}
	return nil
}

func (s *UserService) UpdateUserEmail(ctx context.Context, user *entity.User) error {
	if err := s.repo.UpdateUserEmail(ctx, user.ID, user.Email); err != nil {
		return err
	}
	return nil
}

func (s *UserService) UpdateUserAdminStatus(ctx context.Context, user *entity.User) error {
	if err := s.repo.UpdateUserAdminStatus(ctx, user.ID, user.IsAdmin); err != nil {
		return err
	}
	return nil
}
