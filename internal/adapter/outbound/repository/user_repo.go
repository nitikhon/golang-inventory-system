package repository

import (
	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
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
	err := r.db.Create(user).Error
	if err != nil {
		return &entity.User{}, err
	}
	return user, nil
}

// UpdateUser updates an existing user's details in the database.
func (r *UserRepository) UpdateUser(user *entity.User) (*entity.User, error) {
	result := r.db.Model(&entity.User{}).Where("id = ?", user.ID).Select("*").Updates(user)
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
	err := r.db.Unscoped().Find(&users).Error
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