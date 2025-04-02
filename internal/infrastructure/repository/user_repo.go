package repository

import (
	"github.com/nitikhon/golang-inventory-system/internal/domain"
	"github.com/nitikhon/golang-inventory-system/internal/domain/validation"
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
func (r *UserRepository) CreateUser(user *domain.User) (*domain.User, error) {
	// Validate user fields
	user, err := validation.ValidateAndNormalizeUser(user)
	if err != nil {
		return &domain.User{}, err
	}

	// Check if username already exists
	existingUser, err := r.GetUserByUsername(user.Username)
	if err == nil && existingUser != nil {
		return &domain.User{}, gorm.ErrDuplicatedKey
	}

	// Check if email already exists
	existingUser, err = r.GetUserByEmail(user.Email)
	if err == nil && existingUser != nil {
		return &domain.User{}, gorm.ErrDuplicatedKey
	}

	// Check if phone already exists
	existingUser, err = r.GetUserByPhone(user.Phone)
	if err == nil && existingUser != nil {
		return &domain.User{}, gorm.ErrDuplicatedKey
	}

	// Save user to database
	err = r.db.Create(user).Error
	if err != nil {
		return &domain.User{}, err
	}
	return user, nil
}

// UpdateUser updates an existing user's details in the database.
func (r *UserRepository) UpdateUser(user *domain.User) (*domain.User, error) {
	// Validate user fields
	_, err := validation.ValidateAndNormalizeUser(user)
	if err != nil {
		return &domain.User{}, err
	}

	result := r.db.Model(&domain.User{}).Where("id = ?", user.ID).Updates(user)
	if result.Error != nil {
		return &domain.User{}, result.Error
	}
	if result.RowsAffected == 0 {
		return &domain.User{}, gorm.ErrRecordNotFound
	}
	return user, nil
}

// DeleteUser removes a user from the database by their ID. (soft delete)
func (r *UserRepository) DeleteUser(id int) error {
	result := r.db.Delete(&domain.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetAllUsers retrieves all users from the database.
func (r *UserRepository) GetAllUsers() ([]*domain.User, error) {
	var users []*domain.User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserByID retrieves a user by their ID.
func (r *UserRepository) GetUserByID(id int) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("id = ?", id).Take(&user).Error
	if err != nil {
		return &domain.User{}, err
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by their username.
func (r *UserRepository) GetUserByUsername(username string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("username = ?", username).Take(&user).Error
	if err != nil {
		return &domain.User{}, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by their email.
func (r *UserRepository) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).Take(&user).Error
	if err != nil {
		return &domain.User{}, err
	}
	return &user, nil
}

// GetUserByPhone retrieves a user by their phone number.
func (r *UserRepository) GetUserByPhone(phone string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("phone = ?", phone).Take(&user).Error
	if err != nil {
		return &domain.User{}, err
	}
	return &user, nil
}