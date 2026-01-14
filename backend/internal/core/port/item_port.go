package port

import (
	"context"

	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"gorm.io/gorm"
)

type ItemRepository interface {
	GetAllItems(ctx context.Context, page, limit int, search string) (*entity.PaginationResult[entity.Item], error)
	GetItemByID(ctx context.Context, id uint) (*entity.Item, error)
	GetItemByName(ctx context.Context, name string) (*entity.Item, error)
	Create(ctx context.Context, item *entity.Item) (*entity.Item, error)
	Update(ctx context.Context, item *entity.Item) (*entity.Item, error)
	Delete(ctx context.Context, id uint) error
	GetDB() *gorm.DB
	GetItemByIDForUpdate(tx *gorm.DB, id uint) (*entity.Item, error)
	UpdateWithTx(tx *gorm.DB, item *entity.Item) (*entity.Item, error)
}
