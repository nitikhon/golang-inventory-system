package port

import (
	"context"

	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	UpdateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	DeleteUser(ctx context.Context, id uint) error
	GetAllUsers(ctx context.Context) ([]*entity.User, error)
	GetUserByID(ctx context.Context, id uint) (*entity.User, error)
	GetUserByUsername(ctx context.Context, username string) (*entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (*entity.User, error)
	GetUserByPhone(ctx context.Context, phone string) (*entity.User, error)
	UpdateRefreshToken(ctx context.Context, userID uint, refreshToken string) error
	UpdateUserProfile(ctx context.Context, userID uint, updates map[string]interface{}) (*entity.User, error)
	UpdateUserPassword(ctx context.Context, userID uint, hashedPassword string) error
	UpdateUserEmail(ctx context.Context, userID uint, email string) error
	UpdateUserAdminStatus(ctx context.Context, userID uint, isAdmin bool) error
}
