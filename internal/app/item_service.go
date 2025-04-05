package app

import (
	"github.com/nitikhon/golang-inventory-system/internal/domain"
)

// ItemService provides the use cases for the item domain.
type ItemService struct {
	repo domain.ItemRepository
}

// NewItemService creates a new ItemService instance.
func NewItemService(repo domain.ItemRepository) *ItemService {
	return &ItemService{repo: repo}
}

// FindAll returns all items.
func (s *ItemService) GetAllItems() ([]domain.Item, error) {
	return s.repo.GetAllItems()
}

// FindByID returns an item by its ID.
func (s *ItemService) GetItemByID(id int) (domain.Item, error) {
	return s.repo.GetItemByID(id)
}

// Create creates a new item.
func (s *ItemService) Create(item domain.Item) (domain.Item, error) {
	return s.repo.Create(item)
}

// Update updates an existing item.
func (s *ItemService) Update(item domain.Item) (domain.Item, error) {
	return s.repo.Update(item)
}

// Delete deletes an item by its ID.
func (s *ItemService) Delete(id int) error {
	return s.repo.Delete(id)
}