package repository

import (
	"context"
	"math"

	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"github.com/nitikhon/golang-inventory-system/internal/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
func (r *ItemRepository) GetAllItems(ctx context.Context, page, limit int, search string) (*entity.PaginationResult[entity.Item], error) {
	var items []entity.Item
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Item{})

	if search != "" {
		term := "%" + search + "%"
		query = query.Where("name ILIKE ? OR description ILIKE ?", term, term)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := util.GetOffset(page, limit)
	if err := query.Offset(offset).Limit(limit).Find(&items).Error; err != nil {
		return nil, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return &entity.PaginationResult[entity.Item]{
		Data:       items,
		TotalItems: total,
		TotalPages: totalPages,
		Page:       page,
		Limit:      limit,
	}, nil
}

// FindByID retrieves an item by its ID.
func (r *ItemRepository) GetItemByID(ctx context.Context, id uint) (*entity.Item, error) {
	var item entity.Item
	// Query item by ID
	err := r.db.Where("id = ?", id).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ItemRepository) GetItemByName(ctx context.Context, name string) (*entity.Item, error) {
	var item entity.Item
	err := r.db.WithContext(ctx).Where("name = ?", name).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Create adds a new item to the database.
func (r *ItemRepository) Create(ctx context.Context, item *entity.Item) (*entity.Item, error) {
	// Save item to database
	err := r.db.WithContext(ctx).Create(&item).Error
	if err != nil {
		return nil, err
	}
	return item, nil
}

// Update modifies an existing item in the database.
func (r *ItemRepository) Update(ctx context.Context, item *entity.Item) (*entity.Item, error) {
	// Update item fields
	result := r.db.WithContext(ctx).Model(&entity.Item{}).Where("id = ?", item.ID).Updates(item)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// Retrieve updated item
	var updatedItem entity.Item
	if err := r.db.WithContext(ctx).First(&updatedItem, item.ID).Error; err != nil {
		return nil, err
	}
	return &updatedItem, nil
}

// Delete removes an item by its ID from the database.
func (r *ItemRepository) Delete(ctx context.Context, id uint) error {
	// Delete item by ID
	result := r.db.WithContext(ctx).Delete(&entity.Item{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// GetDB helps other layers to access the db
func (r *ItemRepository) GetDB() *gorm.DB {
	return r.db
}

// GetItemByIDForUpdate retrieves an item by its ID with a row-level lock for update within the given transaction.
// This ensures the selected row is locked for the duration of the transaction to prevent concurrent modifications.
func (r *ItemRepository) GetItemByIDForUpdate(tx *gorm.DB, id uint) (*entity.Item, error) {
	var item entity.Item
	// Acquires a row-level lock on the item with the specified ID using the "FOR UPDATE" clause to prevent concurrent updates during the transaction.
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Where("id = ?", id).Take(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// UpdateWithTx updates an existing Item entity in the database using the provided transaction.
// The function first attempts to update the item, then retrieves and returns the updated record.
func (r *ItemRepository) UpdateWithTx(tx *gorm.DB, item *entity.Item) (*entity.Item, error) {
	result := tx.Model(&entity.Item{}).Where("id = ?", item.ID).Updates(item)
	if result.Error != nil {
		return nil, result.Error
	}
	if result.RowsAffected == 0 {
		return nil, gorm.ErrRecordNotFound
	}

	// Retrieve updated item
	var updatedItem entity.Item
	if err := tx.First(&updatedItem, item.ID).Error; err != nil {
		return nil, err
	}
	return &updatedItem, nil
}
