package repository

import (
	"github.com/nitikhon/golang-inventory-system/internal/domain"
	"gorm.io/gorm"
)

// ItemRepository provides the persistence methods for the item domain.
type ItemRepository struct {
	db *gorm.DB
}

// NewItemRepository creates a new ItemRepository instance.
func NewItemRepository(db *gorm.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

// FindAll returns all items.
func (r *ItemRepository) FindAll() ([]domain.Item, error) {
	var items []domain.Item
	err := r.db.Find(&items).Error
	if err != nil {
		return nil, err
	}

	return items, nil
}

// FindByID returns an item by its ID.
func (r *ItemRepository) FindByID(id int) (domain.Item, error) {
	var item domain.Item
	err := r.db.First(&item, id).Error
	if err != nil {
		return domain.Item{}, err
	}

	return item, nil
}

// Create creates a new item.
func (r *ItemRepository) Create(item domain.Item) (domain.Item, error) {
	err := r.db.Create(&item).Error
	if err != nil {
		return domain.Item{}, err
	}

	return item, nil
}

// Update updates an existing item.
func (r *ItemRepository) Update(item domain.Item) (domain.Item, error) {
	err := r.db.Save(&item).Error
	if err != nil {
		return domain.Item{}, err
	}

	return item, nil
}

// Delete deletes an item by its ID.
func (r *ItemRepository) Delete(id int) error {
	err := r.db.Delete(&domain.Item{}, id).Error
	if err != nil {
		return err
	}

	return nil
}
