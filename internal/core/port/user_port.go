package port

import (
	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
)

type UserRepository interface {
	CreateUser(user *entity.User) (*entity.User, error)
	UpdateUser(user *entity.User) (*entity.User, error)
	DeleteUser(id int) error
	GetAllUsers() ([]*entity.User, error)
	GetUserByID(id int) (*entity.User, error)
	GetUserByUsername(username string) (*entity.User, error)
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByPhone(phone string) (*entity.User, error)
}
