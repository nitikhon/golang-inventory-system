package repository

import (
	"context"
	"errors"

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
func (r *UserRepository) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		return &entity.User{}, err
	}
	return user, nil
}

// UpdateUser updates an existing user's details in the database.
func (r *UserRepository) UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	result := r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", user.ID).
		Select("first_name", "last_name", "phone", "email", "password", "updated_at").
		Updates(user)
	if result.Error != nil {
		return &entity.User{}, result.Error
	}
	if result.RowsAffected == 0 {
		return &entity.User{}, gorm.ErrRecordNotFound
	}
	var updatedUser entity.User
	if err := r.db.WithContext(ctx).First(&updatedUser, user.ID).Error; err != nil {
		return &entity.User{}, err
	}
	return &updatedUser, nil
}

// DeleteUser removes a user from the database by their ID. (soft delete)
func (r *UserRepository) DeleteUser(ctx context.Context, id uint) error {
	result := r.db.WithContext(ctx).Delete(&entity.User{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetAllUsers retrieves all users from the database.
func (r *UserRepository) GetAllUsers(ctx context.Context) ([]*entity.User, error) {
	var users []*entity.User
	err := r.db.WithContext(ctx).Unscoped().Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

// GetUserByID retrieves a user by their ID.
func (r *UserRepository) GetUserByID(ctx context.Context, id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("id = ?", id).Take(&user).Error
	if err != nil {
		return &entity.User{}, err
	}
	return &user, nil
}

// GetUserByUsername retrieves a user by their username.
func (r *UserRepository) GetUserByUsername(ctx context.Context, username string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("username = ?", username).Take(&user).Error
	if err != nil {
		return &entity.User{}, err
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by their email.
func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("email = ?", email).Take(&user).Error
	if err != nil {
		return &entity.User{}, err
	}
	return &user, nil
}

// GetUserByPhone retrieves a user by their phone number.
func (r *UserRepository) GetUserByPhone(ctx context.Context, phone string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Where("phone = ?", phone).Take(&user).Error
	if err != nil {
		return &entity.User{}, err
	}
	return &user, nil
}

// UpdateRefreshToken updates only the refresh token
func (r *UserRepository) UpdateRefreshToken(ctx context.Context, userID uint, refreshToken string) error {
	result := r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", userID).
		Update("refresh_token", refreshToken)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// UpdateUserProfile updates general profile fields
func (r *UserRepository) UpdateUserProfile(ctx context.Context, userID uint, updates map[string]any) (*entity.User, error) {
	allowedFields := map[string]bool{
		"first_name": true,
		"last_name":  true,
		"phone":      true,
	}

	filteredUpdates := make(map[string]any)
	for key, value := range updates {
		if allowedFields[key] {
			filteredUpdates[key] = value
		}
	}

	if len(filteredUpdates) == 0 {
		return nil, errors.New("no valid fields to update")
	}

	result := r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", userID).Updates(filteredUpdates)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	var user entity.User
	err := r.db.WithContext(ctx).First(&user, userID).Error
	return &user, err
}

// UpdateUserPassword updates user password with validation
func (r *UserRepository) UpdateUserPassword(ctx context.Context, userID uint, hashedPassword string) error {
	result := r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", userID).
		Update("password", hashedPassword)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// UpdateUserEmail updates user email with validation
func (r *UserRepository) UpdateUserEmail(ctx context.Context, userID uint, email string) error {
	result := r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", userID).
		Update("email", email)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// UpdateUserAdminStatus updates admin status (admin only)
func (r *UserRepository) UpdateUserAdminStatus(ctx context.Context, userID uint, isAdmin bool) error {
	result := r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", userID).
		Update("is_admin", isAdmin)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
