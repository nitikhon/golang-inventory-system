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
	mockItemRepo.EXPECT().GetAllItems().Return([]*entity.Item{}, nil)

	// act
	items, err := mockItemService.GetAllItems()

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
			mockReturn: &entity.Item{Model: gorm.Model{ID: 1}},
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

func TestCreate(t *testing.T) {
	tests := []struct {
		name            string
		itemName        string
		AvailableAmount int
		mockReturn      *entity.Item
		mockErr         error
		expectErr       bool
	}{
		{
			name:            "success: create an item",
			itemName:        "test",
			AvailableAmount: 10,
			mockReturn:      &entity.Item{Name: "test", AvailableAmount: 10, TotalAmount: 10},
			mockErr:         nil,
			expectErr:       false,
		},
		{
			name:       "failed: empty item name",
			itemName:   "",
			mockReturn: &entity.Item{},
			mockErr:    errors.New("name is required"),
			expectErr:  true,
		},
		{
			name:            "failed: negative AvailableAmount",
			itemName:        "test",
			AvailableAmount: -1,
			mockReturn:      &entity.Item{},
			mockErr:         errors.New("AvailableAmount cannot be negative"),
			expectErr:       true,
		},
		{
			name:            "failed: zero AvailableAmount",
			itemName:        "test",
			AvailableAmount: 0,
			mockReturn:      &entity.Item{},
			mockErr:         errors.New("AvailableAmount cannot be zero"),
			expectErr:       true,
		},
	}

	mockItemService, mockItemRepo := setupItemServiceMock(t)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// arrange
			if !tt.expectErr {
				mockItemRepo.
					EXPECT().
					Create(gomock.Any()).
					Return(tt.mockReturn, tt.mockErr)
			}

			// act
			item, err := mockItemService.Create(&entity.Item{
				Name:            tt.itemName,
				AvailableAmount: tt.AvailableAmount})

			// assert
			if tt.expectErr {
				assert.NotNil(t, err, "expect an error, got nil")
			} else {
				assert.Nil(t, err, "expect nil, got %v", err)
			}

			assert.Equal(t, tt.mockReturn, item, "expect %v, got %v", tt.mockReturn, item)
		})
	}
}

func TestUpdate(t *testing.T) {
	// arrange
	mockItemService, mockItemRepo := setupItemServiceMock(t)
	mockItemRepo.EXPECT().Update(gomock.Any()).Return(&entity.Item{Name: "updatedName"}, nil)

	// act
	item, err := mockItemService.Update(&entity.Item{Name: "updatedName"})

	// assert
	assert.NotNil(t, item, "expect an item, got nil")
	assert.Nil(t, err, fmt.Sprintf("expect an error to be nil, got %v", err))
	assert.Equal(t, "updatedName", item.Name)
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
		Return(&entity.Item{Model: gorm.Model{ID: 1}}, nil)
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
