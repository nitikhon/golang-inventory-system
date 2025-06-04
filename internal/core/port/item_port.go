package port

import (
	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
)

type ItemRepository interface {
	GetAllItems() ([]*entity.Item, error)
	GetItemByID(id uint) (*entity.Item, error)
	Create(item *entity.Item) (*entity.Item, error)
	Update(item *entity.Item) (*entity.Item, error)
	Delete(id int) error
}