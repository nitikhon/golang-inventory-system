package service

import (
	"context"
	"errors"

	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"github.com/nitikhon/golang-inventory-system/internal/core/port"
	"github.com/nitikhon/golang-inventory-system/internal/util/errormap"
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
	GetAllItems(ctx context.Context, page, limit int, search string) (*entity.PaginationResult[entity.Item], error)
	GetItemByID(ctx context.Context, id uint) (*entity.Item, error)
	Create(ctx context.Context, item *entity.Item) (*entity.Item, error)
	Update(ctx context.Context, item *entity.Item) (*entity.Item, error)
	Delete(ctx context.Context, id uint) error
	GetItemByIDForUpdate(tx *gorm.DB, id uint) (*entity.Item, error)
	UpdateWithTx(tx *gorm.DB, item *entity.Item) (*entity.Item, error)
}

// NewItemService creates a new ItemService instance.
func NewItemService(repo port.ItemRepository) *ItemService {
	return &ItemService{repo: repo}
}

// GetAllItems returns all items.
func (s *ItemService) GetAllItems(ctx context.Context, page, limit int, search string) (*entity.PaginationResult[entity.Item], error) {
	return s.repo.GetAllItems(ctx, page, limit, search)
}

// GetItemByID returns an item by its ID.
func (s *ItemService) GetItemByID(ctx context.Context, id uint) (*entity.Item, error) {
	return s.repo.GetItemByID(ctx, id)
}

// Create creates a new item.
func (s *ItemService) Create(ctx context.Context, item *entity.Item) (*entity.Item, error) {
	existingItem, err := s.repo.GetItemByName(ctx, item.Name)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if existingItem != nil {
		return nil, errors.New(errormap.ErrItemNameAlreadyExists)
	}

	return s.repo.Create(ctx, item)
}

// Update updates an existing item.
func (s *ItemService) Update(ctx context.Context, item *entity.Item) (*entity.Item, error) {
	currentItem, err := s.repo.GetItemByID(ctx, item.ID)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	if currentItem == nil {
		return nil, errors.New(errormap.ErrItemNotFound)
	}

	existingItem, err := s.repo.GetItemByName(ctx, item.Name)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if existingItem != nil && existingItem.ID != item.ID {
		return nil, errors.New(errormap.ErrItemNameAlreadyExists)
	}

	return s.repo.Update(ctx, item)
}

// Delete deletes an item by its ID.
func (s *ItemService) Delete(ctx context.Context, id uint) error {
	return s.repo.Delete(ctx, id)
}

// GetItemByIDForUpdate create a transaction to update to handle race condition
func (s *ItemService) GetItemByIDForUpdate(tx *gorm.DB, id uint) (*entity.Item, error) {
	return s.repo.GetItemByIDForUpdate(tx, id)
}

func (s *ItemService) UpdateWithTx(tx *gorm.DB, item *entity.Item) (*entity.Item, error) {
	return s.repo.UpdateWithTx(tx, item)
}
