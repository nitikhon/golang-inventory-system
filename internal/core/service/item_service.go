package service

import (
	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"github.com/nitikhon/golang-inventory-system/internal/core/port"
)

// ItemService provides the use cases for the item entity.
type ItemService struct {
	repo port.ItemRepository
}

// NewItemService creates a new ItemService instance.
func NewItemService(repo port.ItemRepository) *ItemService {
	return &ItemService{repo: repo}
}

// FindAll returns all items.
func (s *ItemService) GetAllItems() ([]entity.Item, error) {
	return s.repo.GetAllItems()
}

// FindByID returns an item by its ID.
func (s *ItemService) GetItemByID(id int) (entity.Item, error) {
	return s.repo.GetItemByID(id)
}

// Create creates a new item.
func (s *ItemService) Create(item entity.Item) (entity.Item, error) {
	return s.repo.Create(item)
}

// Update updates an existing item.
func (s *ItemService) Update(item entity.Item) (entity.Item, error) {
	return s.repo.Update(item)
}

// Delete deletes an item by its ID.
func (s *ItemService) Delete(id int) error {
	return s.repo.Delete(id)
}