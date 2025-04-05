package service

import (
	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"github.com/nitikhon/golang-inventory-system/internal/core/port"
)

// UserService provides the use cases for the user entity.
type UserService struct {
	repo port.UserRepository
}

// NewUserService creates a new UserService instance.
func NewUserService(repo port.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Create creates a new user.
func (s *UserService) Create(user entity.User) (*entity.User, error) {
	return s.repo.CreateUser(&user)
}

// Update updates an existing user.
func (s *UserService) Update(user entity.User) (*entity.User, error) {
	return s.repo.UpdateUser(&user)
}

// Delete deletes a user by their ID.
func (s *UserService) Delete(id int) error {
	return s.repo.DeleteUser(id)
}

// GetAllUsers returns all users.
func (s *UserService) GetAllUsers() ([]*entity.User, error) {
	return s.repo.GetAllUsers()
}

// GetUserByID returns a user by their ID.
func (s *UserService) GetUserByID(id int) (*entity.User, error) {
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