package repository

import (
	"errors"

	"github.com/nitikhon/golang-inventory-system/internal/domain"
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
func (r *ItemRepository) FindAll() ([]domain.Item, error) {
	var items []domain.Item
	// Query all items
	err := r.db.Find(&items).Error
	if err != nil {
		return nil, err
	}
	return items, nil
}

// FindByID retrieves an item by its ID.
func (r *ItemRepository) FindByID(id int) (domain.Item, error) {
	var item domain.Item
	// Query item by ID
	err := r.db.Where("id = ?", id).Take(&item).Error
	if err != nil {
		return domain.Item{}, err
	}
	return item, nil
}

// Create adds a new item to the database.
func (r *ItemRepository) Create(item domain.Item) (domain.Item, error) {
	// Validate item fields
	if item.Name == "" {
		return domain.Item{}, errors.New("name is required")
	}
	if item.Quantity < 0 {
		return domain.Item{}, errors.New("quantity cannot be negative")
	}
	if item.Quantity == 0 {
		return domain.Item{}, errors.New("quantity cannot be zero")
	}

	// Save item to database
	err := r.db.Create(&item).Error
	if err != nil {
		return domain.Item{}, err
	}
	return item, nil
}

// Update modifies an existing item in the database.
func (r *ItemRepository) Update(item domain.Item) (domain.Item, error) {
	// Update item fields
	result := r.db.Model(&domain.Item{}).Where("id = ?", item.ID).Updates(item)
	if result.Error != nil {
		return domain.Item{}, result.Error
	}
	if result.RowsAffected == 0 {
		return domain.Item{}, gorm.ErrRecordNotFound
	}

	// Retrieve updated item
	var updatedItem domain.Item
	if err := r.db.First(&updatedItem, item.ID).Error; err != nil {
		return domain.Item{}, err
	}
	return updatedItem, nil
}

// Delete removes an item by its ID from the database.
func (r *ItemRepository) Delete(id int) error {
	// Delete item by ID
	result := r.db.Delete(&domain.Item{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
