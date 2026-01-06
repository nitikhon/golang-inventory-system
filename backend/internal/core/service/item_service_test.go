package service

import (
	"errors"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	mock_port "github.com/nitikhon/golang-inventory-system/internal/core/port/mock"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func setupItemServiceMock(t *testing.T) (*ItemService, *mock_port.MockItemRepository) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockItemRepo := mock_port.NewMockItemRepository(ctrl)
	service := NewItemService(mockItemRepo)
	return service, mockItemRepo
}

func TestNewItemService(t *testing.T) {
	// arrange
	mockItemService, _ := setupItemServiceMock(t)

	// assert
	assert.NotNil(t, mockItemService)
}

func TestGetAllItems(t *testing.T) {
	// arrange
	mockItemService, mockItemRepo := setupItemServiceMock(t)
	page := 1
	limit := 12
	search := ""
	mockItemRepo.EXPECT().GetAllItems(page, limit, search).Return(&entity.PaginationResult[entity.Item]{}, nil)

	// act
	items, err := mockItemService.GetAllItems(page, limit, search)

	// assert
	assert.NotNil(t, items, fmt.Sprintf("expect empty array, got %v", items))
	assert.Nil(t, err, fmt.Sprintf("expect nil, got %v", err))
}

func TestGetItemById(t *testing.T) {
	tests := []struct {
		name       string
		id         uint
		mockReturn *entity.Item
		mockErr    error
		expectErr  bool
	}{
		{
			name:       "success: found item",
			id:         1,
			mockReturn: &entity.Item{GormModel: entity.GormModel{ID: 1}},
			mockErr:    nil,
			expectErr:  false,
		},
		{
			name:       "fail: not found",
			id:         999,
			mockReturn: &entity.Item{},
			mockErr:    fmt.Errorf("item not found"),
			expectErr:  true,
		},
		{
			name:       "fail: db error",
			id:         2,
			mockReturn: &entity.Item{},
			mockErr:    fmt.Errorf("db error"),
			expectErr:  true,
		},
	}

	mockItemService, mockItemRepo := setupItemServiceMock(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			mockItemRepo.
				EXPECT().
				GetItemByID(tt.id).
				Return(tt.mockReturn, tt.mockErr)

			// act
			item, err := mockItemService.GetItemByID(tt.id)

			// assert
			if tt.expectErr {
				assert.NotNil(t, err, "expect error to be not nil")
			} else {
				assert.Nil(t, err, fmt.Sprintf("expect error to be nil, got %v", err))
			}

			assert.Equal(t, tt.mockReturn, item, "expect an item to be %v, got %v", tt.mockReturn, item)
		})
	}
}

func TestCreateItem_Success(t *testing.T) {
	// arrange
	mockItemService, mockItemRepo := setupItemServiceMock(t)

	inputItem := &entity.Item{
		Name:            "test item",
		Description:     "Test Description",
		AvailableAmount: 10,
		TotalAmount:     15,
		Status:          "available",
	}

	expectedItem := &entity.Item{
		Name:            "test item",
		Description:     "Test Description",
		AvailableAmount: 10,
		TotalAmount:     15,
		Status:          "available",
	}

	mockItemRepo.
		EXPECT().
		GetItemByName("test item").
		Return(nil, nil) // no existing item found

	mockItemRepo.
		EXPECT().
		Create(gomock.Any()).
		Return(expectedItem, nil)

	// act
	result, err := mockItemService.Create(inputItem)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedItem, result)
}

func TestCreateItem_FailureItemAlreadyExists(t *testing.T) {
	// arrange
	mockItemService, mockItemRepo := setupItemServiceMock(t)

	inputItem := &entity.Item{
		Name:            "existing item",
		Description:     "Test Description",
		AvailableAmount: 10,
		TotalAmount:     15,
	}

	existingItem := &entity.Item{
		Name:            "existing item",
		Description:     "Already exists",
		AvailableAmount: 5,
		TotalAmount:     10,
	}

	mockItemRepo.
		EXPECT().
		GetItemByName("existing item"). // normalized name
		Return(existingItem, nil)       // existing item found

	// act
	result, err := mockItemService.Create(inputItem)

	// assert
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "item with this name already exists")
}

func TestCreateItem_FailureGetItemByNameError(t *testing.T) {
	// arrange
	mockItemService, mockItemRepo := setupItemServiceMock(t)

	inputItem := &entity.Item{
		Name:            "test item",
		Description:     "Test Description",
		AvailableAmount: 10,
		TotalAmount:     15,
	}

	mockItemRepo.
		EXPECT().
		GetItemByName("test item"). // normalized name
		Return(nil, errors.New("database connection error"))

	// act
	result, err := mockItemService.Create(inputItem)

	// assert
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database connection error")
}

func TestCreateItem_FailureCreateError(t *testing.T) {
	// arrange
	mockItemService, mockItemRepo := setupItemServiceMock(t)

	inputItem := &entity.Item{
		Name:            "test item",
		Description:     "Test Description",
		AvailableAmount: 10,
		TotalAmount:     15,
	}

	mockItemRepo.
		EXPECT().
		GetItemByName("test item").
		Return(nil, nil) // no existing item found

	mockItemRepo.
		EXPECT().
		Create(gomock.Any()).
		Return(nil, errors.New("failed to insert into database"))

	// act
	result, err := mockItemService.Create(inputItem)

	// assert
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "failed to insert into database")
}

func TestUpdateItem_Success(t *testing.T) {
	// arrange
	mockItemService, mockItemRepo := setupItemServiceMock(t)

	inputItem := &entity.Item{
		GormModel:       entity.GormModel{ID: 1},
		Name:            "Updated Item",
		Description:     "Updated Description",
		AvailableAmount: 8,
		TotalAmount:     12,
		Status:          "available",
	}

	currentItem := &entity.Item{
		GormModel:       entity.GormModel{ID: 1},
		Name:            "old item",
		Description:     "Old Description",
		AvailableAmount: 5,
		TotalAmount:     10,
		Status:          "available",
	}

	expectedItem := &entity.Item{
		GormModel:       entity.GormModel{ID: 1},
		Name:            "updated Item", // service doesn't normalize name yet
		Description:     "Updated Description",
		AvailableAmount: 8,
		TotalAmount:     12,
		Status:          "available",
	}

	// Mock expectations - check item exists first, then check name conflicts
	mockItemRepo.
		EXPECT().
		GetItemByID(uint(1)).
		Return(currentItem, nil) // item exists

	mockItemRepo.
		EXPECT().
		GetItemByName("Updated Item"). // original name (service doesn't normalize yet)
		Return(nil, nil)               // no existing item with this name

	mockItemRepo.
		EXPECT().
		Update(gomock.Any()).
		Return(expectedItem, nil)

	// act
	result, err := mockItemService.Update(inputItem)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedItem, result)
}

func TestUpdateItem_FailureItemNotFound(t *testing.T) {
	// arrange
	mockItemService, mockItemRepo := setupItemServiceMock(t)

	inputItem := &entity.Item{
		GormModel:       entity.GormModel{ID: 999},
		Name:            "Test Item",
		Description:     "Test Description",
		AvailableAmount: 8,
		TotalAmount:     12,
	}

	mockItemRepo.
		EXPECT().
		GetItemByID(uint(999)).
		Return(nil, nil) // item not found

	// act
	result, err := mockItemService.Update(inputItem)

	// assert
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "item not found")
}

func TestUpdateItem_SuccessUpdateToSameName(t *testing.T) {
	// arrange
	mockItemService, mockItemRepo := setupItemServiceMock(t)

	inputItem := &entity.Item{
		GormModel:       entity.GormModel{ID: 1},
		Name:            "same item",
		Description:     "Updated Description",
		AvailableAmount: 8,
		TotalAmount:     12,
		Status:          "available",
	}

	// existed item = to be updated item
	existingItem := &entity.Item{
		GormModel:       entity.GormModel{ID: 1},
		Name:            "same item",
		Description:     "Old Description",
		AvailableAmount: 5,
		TotalAmount:     10,
		Status:          "available",
	}

	expectedItem := &entity.Item{
		GormModel:       entity.GormModel{ID: 1},
		Name:            "same item",
		Description:     "Updated Description",
		AvailableAmount: 8,
		TotalAmount:     12,
		Status:          "available",
	}

	mockItemRepo.
		EXPECT().
		GetItemByID(uint(1)).
		Return(existingItem, nil)

	mockItemRepo.
		EXPECT().
		GetItemByName("same item").
		Return(existingItem, nil) // existing item found with same ID

	mockItemRepo.
		EXPECT().
		Update(gomock.Any()).
		Return(expectedItem, nil)

	// act
	result, err := mockItemService.Update(inputItem)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedItem, result)
}

func TestUpdateItem_FailureNameAlreadyExistsDifferentItem(t *testing.T) {
	// arrange
	mockItemService, mockItemRepo := setupItemServiceMock(t)

	inputItem := &entity.Item{
		GormModel:       entity.GormModel{ID: 1},
		Name:            "existing name",
		Description:     "Updated Description",
		AvailableAmount: 8,
		TotalAmount:     12,
	}

	// Different item with same name
	existingItem := &entity.Item{
		GormModel:       entity.GormModel{ID: 2},
		Name:            "existing name",
		Description:     "Another item",
		AvailableAmount: 5,
		TotalAmount:     10,
	}

	mockItemRepo.
		EXPECT().
		GetItemByID(uint(1)).
		Return(inputItem, nil)

	mockItemRepo.
		EXPECT().
		GetItemByName("existing name").
		Return(existingItem, nil) // existing item found with different ID

	// act
	result, err := mockItemService.Update(inputItem)

	// assert
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Equal(t, errors.New("item with this name already exists"), err)
}

func TestUpdateItem_FailureGetItemByIDError(t *testing.T) {
	// arrange
	mockItemService, mockItemRepo := setupItemServiceMock(t)

	inputItem := &entity.Item{
		GormModel:       entity.GormModel{ID: 1},
		Name:            "Test Item",
		Description:     "Test Description",
		AvailableAmount: 8,
		TotalAmount:     12,
	}

	mockItemRepo.
		EXPECT().
		GetItemByID(uint(1)).
		Return(nil, errors.New("database connection failed"))

	// act
	result, err := mockItemService.Update(inputItem)

	// assert
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database connection failed")
}

func TestUpdateItem_FailureGetItemByNameError(t *testing.T) {
	// arrange
	mockItemService, mockItemRepo := setupItemServiceMock(t)

	inputItem := &entity.Item{
		GormModel:       entity.GormModel{ID: 1},
		Name:            "Test Item",
		Description:     "Test Description",
		AvailableAmount: 8,
		TotalAmount:     12,
	}

	currentItem := &entity.Item{
		GormModel:       entity.GormModel{ID: 1},
		Name:            "old name",
		Description:     "Old Description",
		AvailableAmount: 5,
		TotalAmount:     10,
	}

	mockItemRepo.
		EXPECT().
		GetItemByID(uint(1)).
		Return(currentItem, nil) // item exists

	mockItemRepo.
		EXPECT().
		GetItemByName("Test Item").
		Return(nil, errors.New("database timeout error"))

	// act
	result, err := mockItemService.Update(inputItem)

	// assert
	assert.NotNil(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "database timeout error")
}

func TestDelete(t *testing.T) {
	// arrange
	mockItemService, mockItemRepo := setupItemServiceMock(t)
	mockItemRepo.EXPECT().Delete(gomock.Any()).Return(nil)

	// act
	err := mockItemService.Delete(uint(1))

	// assert
	assert.Nil(t, err, fmt.Sprintf("expect an error to be nil, got %v", err))
}

func TestGetItemByIdForUpdate(t *testing.T) {
	// arrange
	mockItemService, mockItemRepo := setupItemServiceMock(t)
	mockItemRepo.
		EXPECT().
		GetItemByIDForUpdate(gomock.Any(), gomock.Any()).
		Return(&entity.Item{GormModel: entity.GormModel{ID: 1}}, nil)
	tx := &gorm.DB{}

	// act
	item, err := mockItemService.GetItemByIDForUpdate(tx, uint(1))

	// assert
	assert.Nil(t, err, fmt.Sprintf("expect an error to be nil, got %v", err))
	assert.NotNil(t, item)
	assert.Equal(t, uint(1), item.ID)
}

func TestUpdateWithTx(t *testing.T) {
	// arrange
	mockItemService, mockItemRepo := setupItemServiceMock(t)
	mockItemRepo.
		EXPECT().
		UpdateWithTx(gomock.Any(), gomock.Any()).
		Return(&entity.Item{Name: "updatedName"}, nil)
	tx := &gorm.DB{}

	// act
	item, err := mockItemService.UpdateWithTx(tx, &entity.Item{Name: "updatedName"})

	// assert
	assert.Nil(t, err, fmt.Sprintf("expect an error to be nil, got %v", err))
	assert.NotNil(t, item)
	assert.Equal(t, "updatedName", item.Name)
}
