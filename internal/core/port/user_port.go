package port

import (
	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
)

type UserRepository interface {
	CreateUser(user *entity.User) (*entity.User, error)
	UpdateUser(user *entity.User) (*entity.User, error)
	DeleteUser(id uint) error
	GetAllUsers() ([]*entity.User, error)
	GetUserByID(id uint) (*entity.User, error)
	GetUserByUsername(username string) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByPhone(phone string) (*entity.User, error)
	UpdateRefreshToken(id uint, refreshToken string) error
	Login(username, password string) (string, string, error)
	RefreshToken(refreshToken string) (entity.Token, error)
}
