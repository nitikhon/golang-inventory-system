package service

import (
	"errors"
	"strings"

	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"github.com/nitikhon/golang-inventory-system/internal/core/port"
	"gorm.io/gorm"
)

// ItemService provides the use cases for the item entity.
type ItemService struct {
	repo port.ItemRepository
}

// Ensure ItemService implements ItemServiceInterface
var _ ItemServiceInterface = (*ItemService)(nil)

// ItemServiceInterface defines the contract for ItemService.
type ItemServiceInterface interface {
	GetAllItems() ([]*entity.Item, error)
	GetItemByID(id uint) (*entity.Item, error)
	Create(item *entity.Item) (*entity.Item, error)
	Update(item *entity.Item) (*entity.Item, error)
	Delete(id uint) error
	GetItemByIDForUpdate(tx *gorm.DB, id uint) (*entity.Item, error)
	UpdateWithTx(tx *gorm.DB, item *entity.Item) (*entity.Item, error)
}

// NewItemService creates a new ItemService instance.
func NewItemService(repo port.ItemRepository) *ItemService {
	return &ItemService{repo: repo}
}

// GetAllItems returns all items.
func (s *ItemService) GetAllItems() ([]*entity.Item, error) {
	return s.repo.GetAllItems()
}

// GetItemByID returns an item by its ID.
func (s *ItemService) GetItemByID(id uint) (*entity.Item, error) {
	return s.repo.GetItemByID(id)
}

// Create creates a new item.
func (s *ItemService) Create(item *entity.Item) (*entity.Item, error) {
	// Validate item fields
	if item.Name == "" {
		return &entity.Item{}, errors.New("name is required")
	}
	if item.AvailableAmount < 0 {
		return &entity.Item{}, errors.New("AvailableAmount cannot be negative")
	}
	if item.AvailableAmount == 0 {
		return &entity.Item{}, errors.New("AvailableAmount cannot be zero")
	}

	// Normalize item name
	item.Name = strings.ToLower(item.Name)
	item.TotalAmount = item.AvailableAmount

	return s.repo.Create(item)
}

// Update updates an existing item.
func (s *ItemService) Update(item *entity.Item) (*entity.Item, error) {
	item.TotalAmount = item.AvailableAmount
	return s.repo.Update(item)
}

// Delete deletes an item by its ID.
func (s *ItemService) Delete(id uint) error {
	return s.repo.Delete(id)
}

// GetItemByIDForUpdate create a transaction to update to handle race condition
func (s *ItemService) GetItemByIDForUpdate(tx *gorm.DB, id uint) (*entity.Item, error) {
	return s.repo.GetItemByIDForUpdate(tx, id)
}

func (s *ItemService) UpdateWithTx(tx *gorm.DB, item *entity.Item) (*entity.Item, error) {
	return s.repo.UpdateWithTx(tx, item)
}
