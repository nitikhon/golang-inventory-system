package app

import (
	"github.com/nitikhon/golang-inventory-system/internal/domain"
)

// UserService provides the use cases for the user domain.
type UserService struct {
	repo domain.UserRepository
}

// NewUserService creates a new UserService instance.
func NewUserService(repo domain.UserRepository) *UserService {
	return &UserService{repo: repo}
}

// Create creates a new user.
func (s *UserService) Create(user domain.User) (*domain.User, error) {
	return s.repo.CreateUser(&user)
}

// Update updates an existing user.
func (s *UserService) Update(user domain.User) (*domain.User, error) {
	return s.repo.UpdateUser(&user)
}

// Delete deletes a user by their ID.
func (s *UserService) Delete(id int) error {
	return s.repo.DeleteUser(id)
}

// GetAllUsers returns all users.
func (s *UserService) GetAllUsers() ([]*domain.User, error) {
	return s.repo.GetAllUsers()
}

// GetUserByID returns a user by their ID.
func (s *UserService) GetUserByID(id int) (*domain.User, error) {
	return s.repo.GetUserByID(id)
}

// GetUserByUsername retrieves a user by their username.
func (s *UserService) GetUserByUsername(username string) (*domain.User, error) {
	return s.repo.GetUserByUsername(username)
}

// GetUserByEmail retrieves a user by their email.
func (s *UserService) GetUserByEmail(email string) (*domain.User, error) {
	return s.repo.GetUserByEmail(email)
}

// GetUserByPhone retrieves a user by their phone number.
func (s *UserService) GetUserByPhone(phone string) (*domain.User, error) {
	return s.repo.GetUserByPhone(phone)
}