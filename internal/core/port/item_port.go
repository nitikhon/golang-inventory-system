package port

import (
	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"gorm.io/gorm"
)

type ItemRepository interface {
	GetAllItems() ([]*entity.Item, error)
	GetItemByID(id uint) (*entity.Item, error)
	Create(item *entity.Item) (*entity.Item, error)
	Update(item *entity.Item) (*entity.Item, error)
	Delete(id uint) error
	GetDB() *gorm.DB
	GetItemByIDForUpdate(tx *gorm.DB, id uint) (*entity.Item, error)
	UpdateWithTx(tx *gorm.DB, item *entity.Item) (*entity.Item, error)
}
