package repository

import (
	"errors"

	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"github.com/nitikhon/golang-inventory-system/internal/util"
	"github.com/nitikhon/golang-inventory-system/internal/util/validation"
	"gorm.io/gorm"
)

// UserRepository provides methods to interact with the database for users.
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new instance of UserRepository.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// CreateUser adds a new user to the database.
func (r *UserRepository) CreateUser(user *entity.User) (*entity.User, error) {
	// Validate user fields
	_, err := validation.ValidateAndNormalizeUser(user)
	if err != nil {
		return &entity.User{}, err
	}

	// Check if username already exists
	existingUser, err := r.GetUserByUsername(user.Username)
	if err == nil && existingUser != nil {
		return &entity.User{}, gorm.ErrDuplicatedKey
	}

	// Check if email already exists
	existingUser, err = r.GetUserByEmail(user.Email)
	if err == nil && existingUser != nil {
		return &entity.User{}, gorm.ErrDuplicatedKey
	}

	// Check if phone already exists
	existingUser, err = r.GetUserByPhone(user.Phone)
	if err == nil && existingUser != nil {
		return &entity.User{}, gorm.ErrDuplicatedKey
	}

	// Hash the password before saving
	hashedPassword, err := util.HashPassword(user.Password)
	if err != nil {
		return &entity.User{}, err
	}
	user.Password = hashedPassword

	// Save user to database
	err = r.db.Create(user).Error
	if err != nil {
		return &entity.User{}, err
	}
	return user, nil
}

// UpdateUser updates an existing user's details in the database.
func (r *UserRepository) UpdateUser(user *entity.User) (*entity.User, error) {
	// Validate user fields
	_, err := validation.ValidateAndNormalizeUser(user)
	if err != nil {
		return &entity.User{}, err
	}

	result := r.db.Model(&entity.User{}).Where("id = ?", user.ID).Updates(user)
	if result.Error != nil {
		return &entity.User{}, result.Error
	}
	if result.RowsAffected == 0 {
		return &entity.User{}, gorm.ErrRecordNotFound
	}
	return user, nil
}

// DeleteUser removes a user from the database by their ID. (soft delete)
func (r *UserRepository) DeleteUser(id uint) error {
	result := r.db.Delete(&entity.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetAllUsers retrieves all users from the database.
func (r *UserRepository) GetAllUsers() ([]*entity.User, error) {
	var users []*entity.User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserByID retrieves a user by their ID.
func (r *UserRepository) GetUserByID(id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("id = ?", id).Take(&user).Error
	if err != nil {
		return &entity.User{}, err
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by their username.
func (r *UserRepository) GetUserByUsername(username string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("username = ?", username).Take(&user).Error
	if err != nil {
		return &entity.User{}, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by their email.
func (r *UserRepository) GetUserByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).Take(&user).Error
	if err != nil {
		return &entity.User{}, err
	}
	return &user, nil
}

// GetUserByPhone retrieves a user by their phone number.
func (r *UserRepository) GetUserByPhone(phone string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("phone = ?", phone).Take(&user).Error
	if err != nil {
		return &entity.User{}, err
	}
	return &user, nil
}

// UpdateRefreshToken updates the refresh token for a user.
func (r *UserRepository) UpdateRefreshToken(userID uint, refreshToken string) error {
	result := r.db.Model(&entity.User{}).Where("id = ?", userID).Update("refresh_token", refreshToken)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Login checks if the user exists and if the password is correct.
func (r *UserRepository) Login(username, password string) (string, string, error) {
	// Check if user exists
	user, err := r.GetUserByUsername(username)
	if err != nil {
		return "", "", err
	}

	if user == nil {
		return "", "", gorm.ErrRecordNotFound
	}

	// Check if password is correct
	if err := util.CheckPasswordHash(user.Password, password); err != nil {
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

	// Save refresh token
	err = r.UpdateRefreshToken(user.ID, refreshToken)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// RefreshToken validates the provided refresh token, retrieves the associated user,
// and generates a new pair of access and refresh tokens.
func (r *UserRepository) RefreshToken(refreshToken string) (entity.Token, error) {
	// Validate the Refresh Token
	userID, err := util.ValidateRefreshToken(refreshToken)
	if err != nil {
		return entity.Token{}, err
	}

	// Retrieve the user
	user, err := r.GetUserByID(userID)
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

	// Save new refresh token
	err = r.UpdateRefreshToken(user.ID, newRefreshToken)
	if err != nil {
		return entity.Token{}, err
	}

	return entity.Token{AccessToken: accessToken, RefreshToken: newRefreshToken}, nil
}

