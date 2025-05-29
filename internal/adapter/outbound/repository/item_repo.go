package repository

import (
	"errors"
	"strings"

	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"gorm.io/gorm"
)

// ItemRepository provides methods to interact with the database for items.
type ItemRepository struct {
	db *gorm.DB
}

// NewItemRepository creates a new instance of ItemRepository.
func NewItemRepository(db *gorm.DB) *ItemRepository {
	return &ItemRepository{db: db}
}

// FindAll retrieves all items from the database.
func (r *ItemRepository) GetAllItems() ([]entity.Item, error) {
	var items []entity.Item
	// Query all items
	err := r.db.Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

// FindByID retrieves an item by its ID.
func (r *ItemRepository) GetItemByID(id int) (entity.Item, error) {
	var item entity.Item
	// Query item by ID
	err := r.db.Where("id = ?", id).Take(&item).Error
	if err != nil {
		return entity.Item{}, err
	}
	return item, nil
}

// Create adds a new item to the database.
func (r *ItemRepository) Create(item entity.Item) (entity.Item, error) {
	// Validate item fields
	if item.Name == "" {
		return entity.Item{}, errors.New("name is required")
	}
	if item.AvailableAmount < 0 {
		return entity.Item{}, errors.New("AvailableAmount cannot be negative")
	}
	if item.AvailableAmount == 0 {
		return entity.Item{}, errors.New("AvailableAmount cannot be zero")
	}

	// Normalize item name
	item.Name = strings.ToLower(item.Name)

	// Save item to database
	err := r.db.Create(&item).Error
	if err != nil {
		return entity.Item{}, err
	}
	return item, nil
}

// Update modifies an existing item in the database.
func (r *ItemRepository) Update(item entity.Item) (entity.Item, error) {
	// Update item fields
	result := r.db.Model(&entity.Item{}).Where("id = ?", item.ID).Updates(item)
	if result.Error != nil {
		return entity.Item{}, result.Error
	}
	if result.RowsAffected == 0 {
		return entity.Item{}, gorm.ErrRecordNotFound
	}

	// Retrieve updated item
	var updatedItem entity.Item
	if err := r.db.First(&updatedItem, item.ID).Error; err != nil {
		return entity.Item{}, err
	}
	return updatedItem, nil
}

// Delete removes an item by its ID from the database.
func (r *ItemRepository) Delete(id int) error {
	// Delete item by ID
	result := r.db.Delete(&entity.Item{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
