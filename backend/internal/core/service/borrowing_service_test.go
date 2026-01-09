package service

import (
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	mock_port "github.com/nitikhon/golang-inventory-system/internal/core/port/mock"
	"github.com/nitikhon/golang-inventory-system/internal/util/errormap"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupBorrowingServiceMock(t *testing.T) (
	*BorrowingService,
	*mock_port.MockBorrowingRepository,
	*mock_port.MockItemRepository,
	*mock_port.MockUserRepository) {

	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	mockItemRepo := mock_port.NewMockItemRepository(ctrl)
	mockUserRepo := mock_port.NewMockUserRepository(ctrl)
	mockBorrowingRepo := mock_port.NewMockBorrowingRepository(ctrl)

	service := NewBorrowingService(mockBorrowingRepo, mockItemRepo, mockUserRepo)
	return service, mockBorrowingRepo, mockItemRepo, mockUserRepo
}

// setupMockDB creates a mock database with sqlmock for testing transactions
func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	return gormDB, mock
}

func TestNewBorrowingService(t *testing.T) {
	service, _, _, _ := setupBorrowingServiceMock(t)

	assert.NotNil(t, service)
}

func TestBorrowItem_Success(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, mockUserRepo := setupBorrowingServiceMock(t)

	// Mock data based on setup.go seed data
	mockUser := &entity.User{
		GormModel: entity.GormModel{ID: 2},
		Username:  "test_user",
		Email:     "user@test.com",
		Phone:     "0987654321",
		FirstName: "Test",
		LastName:  "User",
		IsAdmin:   false,
	}

	mockItem := &entity.Item{
		GormModel:       entity.GormModel{ID: 1},
		Name:            "office chair ergonomic",
		Description:     "Ergonomic office chair with lumbar support and adjustable height",
		AvailableAmount: 4,
		TotalAmount:     6,
		Status:          "available",
	}

	borrowingInput := entity.BorrowRequest{
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
	}

	expectedBorrowing := &entity.Borrowing{
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_PENDING,
	}

	mockUserRepo.EXPECT().GetUserByID(borrowingInput.UserID).Return(mockUser, nil)
	mockItemRepo.EXPECT().GetItemByID(borrowingInput.ItemID).Return(mockItem, nil)
	mockBorrowingRepo.EXPECT().HasActiveBorrowing(borrowingInput.UserID, borrowingInput.ItemID).Return(false, nil)
	mockBorrowingRepo.EXPECT().BorrowItem(gomock.Any()).Return(expectedBorrowing, nil)

	// act
	result, err := service.BorrowItem(borrowingInput)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedBorrowing.UserID, result.UserID)
	assert.Equal(t, expectedBorrowing.ItemID, result.ItemID)
}

func TestBorrowItem_UserNotExist(t *testing.T) {
	// arrange
	service, _, _, mockUserRepo := setupBorrowingServiceMock(t)

	borrowingInput := entity.BorrowRequest{
		UserID:          999,
		ItemID:          1,
		BorrowingAmount: 1,
	}

	mockUserRepo.EXPECT().GetUserByID(borrowingInput.UserID).Return(nil, nil)

	// act
	result, err := service.BorrowItem(borrowingInput)

	// assert
	assert.Nil(t, result)
	assert.EqualError(t, err, errormap.ErrUserNotExist)
}

func TestBorrowItem_UserRepoError(t *testing.T) {
	// arrange
	service, _, _, mockUserRepo := setupBorrowingServiceMock(t)

	borrowingInput := entity.BorrowRequest{
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
	}
	mockErr := errors.New("database error")

	mockUserRepo.EXPECT().GetUserByID(borrowingInput.UserID).Return(nil, mockErr)

	// act
	result, err := service.BorrowItem(borrowingInput)

	// assert
	assert.Nil(t, result)
	assert.EqualError(t, err, mockErr.Error())
}

func TestBorrowItem_ItemNotAvailable(t *testing.T) {
	// arrange
	service, _, mockItemRepo, mockUserRepo := setupBorrowingServiceMock(t)

	mockUser := &entity.User{
		GormModel: entity.GormModel{ID: 2},
		Username:  "test_user",
	}

	// Item with status "borrowed" (based on setup.go - conference camera 4k)
	mockItem := &entity.Item{
		GormModel:       entity.GormModel{ID: 4},
		Name:            "conference camera 4k",
		AvailableAmount: 0,
		TotalAmount:     4,
		Status:          "borrowed",
	}

	borrowingInput := entity.BorrowRequest{
		UserID:          2,
		ItemID:          4,
		BorrowingAmount: 1,
	}

	mockUserRepo.EXPECT().GetUserByID(borrowingInput.UserID).Return(mockUser, nil)
	mockItemRepo.EXPECT().GetItemByID(borrowingInput.ItemID).Return(mockItem, nil)

	// act
	result, err := service.BorrowItem(borrowingInput)

	// assert
	assert.Nil(t, result)
	assert.EqualError(t, err, errormap.ErrItemNotAvailable)
}

func TestBorrowItem_ItemNotExist(t *testing.T) {
	// arrange
	service, _, mockItemRepo, mockUserRepo := setupBorrowingServiceMock(t)

	mockUser := &entity.User{
		GormModel: entity.GormModel{ID: 2},
		Username:  "test_user",
	}

	borrowingInput := entity.BorrowRequest{
		UserID:          2,
		ItemID:          999,
		BorrowingAmount: 1,
	}

	mockUserRepo.EXPECT().GetUserByID(borrowingInput.UserID).Return(mockUser, nil)
	mockItemRepo.EXPECT().GetItemByID(borrowingInput.ItemID).Return(nil, nil)

	// act
	result, err := service.BorrowItem(borrowingInput)

	// assert
	assert.Nil(t, result)
	assert.EqualError(t, err, errormap.ErrItemNotAvailable)
}

func TestBorrowItem_ItemRepoError(t *testing.T) {
	// arrange
	service, _, mockItemRepo, mockUserRepo := setupBorrowingServiceMock(t)

	mockUser := &entity.User{
		GormModel: entity.GormModel{ID: 2},
		Username:  "test_user",
	}

	borrowingInput := entity.BorrowRequest{
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
	}
	mockErr := errors.New("database error")

	mockUserRepo.EXPECT().GetUserByID(borrowingInput.UserID).Return(mockUser, nil)
	mockItemRepo.EXPECT().GetItemByID(borrowingInput.ItemID).Return(nil, mockErr)

	// act
	result, err := service.BorrowItem(borrowingInput)

	// assert
	assert.Nil(t, result)
	assert.EqualError(t, err, mockErr.Error())
}

func TestBorrowItem_ItemNotEnough(t *testing.T) {
	// arrange
	service, _, mockItemRepo, mockUserRepo := setupBorrowingServiceMock(t)

	mockUser := &entity.User{
		GormModel: entity.GormModel{ID: 2},
		Username:  "test_user",
	}

	// Item with limited availability (based on setup.go - portable projector)
	mockItem := &entity.Item{
		GormModel:       entity.GormModel{ID: 6},
		Name:            "portable projector",
		AvailableAmount: 1,
		TotalAmount:     2,
		Status:          "available",
	}

	borrowingInput := entity.BorrowRequest{
		UserID:          2,
		ItemID:          6,
		BorrowingAmount: 5, // Requesting more than available
	}

	mockUserRepo.EXPECT().GetUserByID(borrowingInput.UserID).Return(mockUser, nil)
	mockItemRepo.EXPECT().GetItemByID(borrowingInput.ItemID).Return(mockItem, nil)

	// act
	result, err := service.BorrowItem(borrowingInput)

	// assert
	assert.Nil(t, result)
	assert.EqualError(t, err, errormap.ErrItemNotEnough)
}

func TestBorrowItem_AlreadyBorrowed(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, mockUserRepo := setupBorrowingServiceMock(t)

	mockUser := &entity.User{
		GormModel: entity.GormModel{ID: 2},
		Username:  "test_user",
	}

	mockItem := &entity.Item{
		GormModel:       entity.GormModel{ID: 1},
		Name:            "office chair ergonomic",
		AvailableAmount: 4,
		TotalAmount:     6,
		Status:          "available",
	}

	borrowingInput := entity.BorrowRequest{
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
	}

	mockUserRepo.EXPECT().GetUserByID(borrowingInput.UserID).Return(mockUser, nil)
	mockItemRepo.EXPECT().GetItemByID(borrowingInput.ItemID).Return(mockItem, nil)
	mockBorrowingRepo.EXPECT().HasActiveBorrowing(borrowingInput.UserID, borrowingInput.ItemID).Return(true, nil)

	// act
	result, err := service.BorrowItem(borrowingInput)

	// assert
	assert.Nil(t, result)
	assert.EqualError(t, err, errormap.ErrAlreadyBorrowed)
}

func TestBorrowItem_HasActiveBorrowingRepoError(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, mockUserRepo := setupBorrowingServiceMock(t)

	mockUser := &entity.User{
		GormModel: entity.GormModel{ID: 2},
		Username:  "test_user",
	}

	mockItem := &entity.Item{
		GormModel:       entity.GormModel{ID: 1},
		Name:            "office chair ergonomic",
		AvailableAmount: 4,
		TotalAmount:     6,
		Status:          "available",
	}

	borrowingInput := entity.BorrowRequest{
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
	}
	mockErr := errors.New("database error")

	mockUserRepo.EXPECT().GetUserByID(borrowingInput.UserID).Return(mockUser, nil)
	mockItemRepo.EXPECT().GetItemByID(borrowingInput.ItemID).Return(mockItem, nil)
	mockBorrowingRepo.EXPECT().HasActiveBorrowing(borrowingInput.UserID, borrowingInput.ItemID).Return(false, mockErr)

	// act
	result, err := service.BorrowItem(borrowingInput)

	// assert
	assert.Nil(t, result)
	assert.EqualError(t, err, mockErr.Error())
}

func TestBorrowItem_BorrowRepoError(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, mockUserRepo := setupBorrowingServiceMock(t)

	mockUser := &entity.User{
		GormModel: entity.GormModel{ID: 2},
		Username:  "test_user",
	}

	mockItem := &entity.Item{
		GormModel:       entity.GormModel{ID: 1},
		Name:            "office chair ergonomic",
		AvailableAmount: 4,
		TotalAmount:     6,
		Status:          "available",
	}

	borrowingInput := entity.BorrowRequest{
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
	}
	mockErr := errors.New("database error")

	mockUserRepo.EXPECT().GetUserByID(borrowingInput.UserID).Return(mockUser, nil)
	mockItemRepo.EXPECT().GetItemByID(borrowingInput.ItemID).Return(mockItem, nil)
	mockBorrowingRepo.EXPECT().HasActiveBorrowing(borrowingInput.UserID, borrowingInput.ItemID).Return(false, nil)
	mockBorrowingRepo.EXPECT().BorrowItem(gomock.Any()).Return(nil, mockErr)

	// act
	result, err := service.BorrowItem(borrowingInput)

	// assert
	assert.Nil(t, result)
	assert.EqualError(t, err, mockErr.Error())
}

func TestBorrowItem_CanBorrowDifferentItem(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, mockUserRepo := setupBorrowingServiceMock(t)

	mockUser := &entity.User{
		GormModel: entity.GormModel{ID: 2},
		Username:  "test_user",
	}

	// First item - wireless keyboard (based on setup.go)
	mockItem := &entity.Item{
		GormModel:       entity.GormModel{ID: 3},
		Name:            "wireless keyboard",
		AvailableAmount: 7,
		TotalAmount:     10,
		Status:          "available",
	}

	borrowingInput := entity.BorrowRequest{
		UserID:          2,
		ItemID:          3, // Different item
		BorrowingAmount: 2,
	}

	expectedBorrowing := &entity.Borrowing{
		UserID:          2,
		ItemID:          3,
		BorrowingAmount: 2,
		BorrowingStatus: entity.BORROWING_PENDING,
	}

	mockUserRepo.EXPECT().GetUserByID(borrowingInput.UserID).Return(mockUser, nil)
	mockItemRepo.EXPECT().GetItemByID(borrowingInput.ItemID).Return(mockItem, nil)
	mockBorrowingRepo.EXPECT().HasActiveBorrowing(borrowingInput.UserID, borrowingInput.ItemID).Return(false, nil)
	mockBorrowingRepo.EXPECT().BorrowItem(gomock.Any()).Return(expectedBorrowing, nil)

	// act
	result, err := service.BorrowItem(borrowingInput)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedBorrowing.ItemID, result.ItemID)
}

func TestBorrowItem_InvalidDueDateFormat(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, mockUserRepo := setupBorrowingServiceMock(t)

	mockUser := &entity.User{
		GormModel: entity.GormModel{ID: 2},
		Username:  "test_user",
	}

	mockItem := &entity.Item{
		GormModel:       entity.GormModel{ID: 1},
		Name:            "office chair ergonomic",
		AvailableAmount: 4,
		TotalAmount:     6,
		Status:          "available",
	}

	borrowingInput := entity.BorrowRequest{
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		DueDate:         "2025/01/01",
	}

	mockUserRepo.EXPECT().GetUserByID(borrowingInput.UserID).Return(mockUser, nil)
	mockItemRepo.EXPECT().GetItemByID(borrowingInput.ItemID).Return(mockItem, nil)
	mockBorrowingRepo.EXPECT().HasActiveBorrowing(borrowingInput.UserID, borrowingInput.ItemID).Return(false, nil)

	// act
	result, err := service.BorrowItem(borrowingInput)

	// assert
	assert.Nil(t, result)
	assert.EqualError(t, err, errormap.ErrInvalidDueDateFormat)
}

func TestBorrowItem_InvalidDueDate(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, mockUserRepo := setupBorrowingServiceMock(t)

	mockUser := &entity.User{
		GormModel: entity.GormModel{ID: 2},
		Username:  "test_user",
	}

	mockItem := &entity.Item{
		GormModel:       entity.GormModel{ID: 1},
		Name:            "office chair ergonomic",
		AvailableAmount: 4,
		TotalAmount:     6,
		Status:          "available",
	}

	borrowingInput := entity.BorrowRequest{
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		DueDate:         time.Now().AddDate(0, 0, -7).Format(time.RFC3339),
	}

	mockUserRepo.EXPECT().GetUserByID(borrowingInput.UserID).Return(mockUser, nil)
	mockItemRepo.EXPECT().GetItemByID(borrowingInput.ItemID).Return(mockItem, nil)
	mockBorrowingRepo.EXPECT().HasActiveBorrowing(borrowingInput.UserID, borrowingInput.ItemID).Return(false, nil)

	// act
	result, err := service.BorrowItem(borrowingInput)

	// assert
	assert.Nil(t, result)
	assert.EqualError(t, err, errormap.ErrInvalidDueDateValue)
}

func TestBorrowItem_CanBorrowIfPreviousBorrowingNotPending(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, mockUserRepo := setupBorrowingServiceMock(t)

	mockUser := &entity.User{
		GormModel: entity.GormModel{ID: 2},
		Username:  "test_user",
	}

	mockItem := &entity.Item{
		GormModel:       entity.GormModel{ID: 1},
		Name:            "office chair ergonomic",
		AvailableAmount: 4,
		TotalAmount:     6,
		Status:          "available",
	}

	borrowingInput := entity.BorrowRequest{
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
	}

	expectedBorrowing := &entity.Borrowing{
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_PENDING,
	}

	mockUserRepo.EXPECT().GetUserByID(borrowingInput.UserID).Return(mockUser, nil)
	mockItemRepo.EXPECT().GetItemByID(borrowingInput.ItemID).Return(mockItem, nil)
	mockBorrowingRepo.EXPECT().HasActiveBorrowing(borrowingInput.UserID, borrowingInput.ItemID).Return(false, nil)
	mockBorrowingRepo.EXPECT().BorrowItem(gomock.Any()).Return(expectedBorrowing, nil)

	// act
	result, err := service.BorrowItem(borrowingInput)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

func TestApproveBorrowing_Success(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, mockUserRepo := setupBorrowingServiceMock(t)
	mockDB, sqlMock := setupMockDB(t)

	borrowingId := uint(1)
	approverId := uint(1) // test_admin from setup.go

	mockApprover := &entity.User{
		GormModel: entity.GormModel{ID: 1},
		Username:  "test_admin",
		IsAdmin:   true,
	}

	mockBorrowing := &entity.Borrowing{
		GormModel:       entity.GormModel{ID: 1},
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_PENDING,
	}

	mockItem := &entity.Item{
		GormModel:       entity.GormModel{ID: 1},
		Name:            "office chair ergonomic",
		AvailableAmount: 4,
		TotalAmount:     6,
		Status:          "available",
	}

	expectedResult := &entity.Borrowing{
		GormModel:       entity.GormModel{ID: 1},
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_ACTIVE,
		ApprovedBy:      approverId,
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()

	mockItemRepo.EXPECT().GetDB().Return(mockDB)
	mockBorrowingRepo.EXPECT().GetBorrowingByID(borrowingId).Return(mockBorrowing, nil)
	mockUserRepo.EXPECT().GetUserByID(approverId).Return(mockApprover, nil)
	mockItemRepo.EXPECT().GetItemByIDForUpdate(gomock.Any(), mockBorrowing.ItemID).Return(mockItem, nil)
	mockItemRepo.EXPECT().UpdateWithTx(gomock.Any(), gomock.Any()).Return(mockItem, nil)
	mockBorrowingRepo.EXPECT().ApproveBorrowingWithTx(gomock.Any(), borrowingId, approverId).Return(expectedResult, nil)

	// act
	result, err := service.ApproveBorrowing(borrowingId, approverId)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, approverId, result.ApprovedBy)

	// Verify all sqlmock expectations were met
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestApproveBorrowing_BorrowingNotExist(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, _ := setupBorrowingServiceMock(t)
	mockDB, sqlMock := setupMockDB(t)

	borrowingId := uint(999)
	approverId := uint(1)

	// Setup sqlmock expectations for transaction (Begin will be called, then Rollback in defer)
	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	mockItemRepo.EXPECT().GetDB().Return(mockDB)
	mockBorrowingRepo.EXPECT().GetBorrowingByID(borrowingId).Return(nil, nil)

	// act
	result, err := service.ApproveBorrowing(borrowingId, approverId)

	// assert
	assert.Nil(t, result)
	assert.EqualError(t, err, errormap.ErrBorrowingNotExist)
}

func TestApproveBorrowing_BorrowingRepoError(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, _ := setupBorrowingServiceMock(t)
	mockDB, sqlMock := setupMockDB(t)

	borrowingId := uint(1)
	approverId := uint(1)
	mockErr := errors.New("database error")

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	mockItemRepo.EXPECT().GetDB().Return(mockDB)
	mockBorrowingRepo.EXPECT().GetBorrowingByID(borrowingId).Return(nil, mockErr)

	// act
	result, err := service.ApproveBorrowing(borrowingId, approverId)

	// assert
	assert.Nil(t, result)
	assert.EqualError(t, err, mockErr.Error())
}

func TestApproveBorrowing_ApproverNotExist(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, mockUserRepo := setupBorrowingServiceMock(t)
	mockDB, sqlMock := setupMockDB(t)

	borrowingId := uint(1)
	approverId := uint(999) // Non-existent approver

	mockBorrowing := &entity.Borrowing{
		GormModel:       entity.GormModel{ID: 1},
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_PENDING,
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	mockItemRepo.EXPECT().GetDB().Return(mockDB)
	mockBorrowingRepo.EXPECT().GetBorrowingByID(borrowingId).Return(mockBorrowing, nil)
	mockUserRepo.EXPECT().GetUserByID(approverId).Return(nil, nil)

	// act
	result, err := service.ApproveBorrowing(borrowingId, approverId)

	// assert
	assert.Nil(t, result)
	assert.EqualError(t, err, errormap.ErrApproverNotExist)
}

func TestApproveBorrowing_ApproverRepoError(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, mockUserRepo := setupBorrowingServiceMock(t)
	mockDB, sqlMock := setupMockDB(t)

	borrowingId := uint(1)
	approverId := uint(1)
	mockErr := errors.New("database error")

	mockBorrowing := &entity.Borrowing{
		GormModel:       entity.GormModel{ID: 1},
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_PENDING,
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	mockItemRepo.EXPECT().GetDB().Return(mockDB)
	mockBorrowingRepo.EXPECT().GetBorrowingByID(borrowingId).Return(mockBorrowing, nil)
	mockUserRepo.EXPECT().GetUserByID(approverId).Return(nil, mockErr)

	// act
	result, err := service.ApproveBorrowing(borrowingId, approverId)

	// assert
	assert.Nil(t, result)
	assert.EqualError(t, err, mockErr.Error())
}

func TestApproveBorrowing_BorrowingNotPending_AlreadyApproved(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, mockUserRepo := setupBorrowingServiceMock(t)
	mockDB, sqlMock := setupMockDB(t)

	borrowingId := uint(1)
	approverId := uint(1)

	mockApprover := &entity.User{
		GormModel: entity.GormModel{ID: 1},
		Username:  "test_admin",
		IsAdmin:   true,
	}

	// Borrowing already approved
	mockBorrowing := &entity.Borrowing{
		GormModel:       entity.GormModel{ID: 1},
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_ACTIVE,
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	mockItemRepo.EXPECT().GetDB().Return(mockDB)
	mockBorrowingRepo.EXPECT().GetBorrowingByID(borrowingId).Return(mockBorrowing, nil)
	mockUserRepo.EXPECT().GetUserByID(approverId).Return(mockApprover, nil)

	// act
	result, err := service.ApproveBorrowing(borrowingId, approverId)

	// assert
	assert.Nil(t, result)
	assert.EqualError(t, err, errormap.ErrBorrowingNotPending)
}

func TestApproveBorrowing_BorrowingNotPending_AlreadyRejected(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, mockUserRepo := setupBorrowingServiceMock(t)
	mockDB, sqlMock := setupMockDB(t)

	borrowingId := uint(1)
	approverId := uint(1)

	mockApprover := &entity.User{
		GormModel: entity.GormModel{ID: 1},
		Username:  "test_admin",
		IsAdmin:   true,
	}

	// Borrowing already rejected
	mockBorrowing := &entity.Borrowing{
		GormModel:       entity.GormModel{ID: 1},
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_CANCELLED,
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	mockItemRepo.EXPECT().GetDB().Return(mockDB)
	mockBorrowingRepo.EXPECT().GetBorrowingByID(borrowingId).Return(mockBorrowing, nil)
	mockUserRepo.EXPECT().GetUserByID(approverId).Return(mockApprover, nil)

	// act
	result, err := service.ApproveBorrowing(borrowingId, approverId)

	// assert
	assert.Nil(t, result)
	assert.EqualError(t, err, errormap.ErrBorrowingNotPending)
}

func TestApproveBorrowing_ItemNotExist(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, mockUserRepo := setupBorrowingServiceMock(t)
	mockDB, sqlMock := setupMockDB(t)

	borrowingId := uint(1)
	approverId := uint(1)

	mockApprover := &entity.User{
		GormModel: entity.GormModel{ID: 1},
		Username:  "test_admin",
		IsAdmin:   true,
	}

	mockBorrowing := &entity.Borrowing{
		GormModel:       entity.GormModel{ID: 1},
		UserID:          2,
		ItemID:          999, // Non-existent item
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_PENDING,
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	mockItemRepo.EXPECT().GetDB().Return(mockDB)
	mockBorrowingRepo.EXPECT().GetBorrowingByID(borrowingId).Return(mockBorrowing, nil)
	mockUserRepo.EXPECT().GetUserByID(approverId).Return(mockApprover, nil)
	mockItemRepo.EXPECT().GetItemByIDForUpdate(gomock.Any(), mockBorrowing.ItemID).Return(nil, nil)

	// act
	result, err := service.ApproveBorrowing(borrowingId, approverId)

	// assert
	assert.Nil(t, result)
	assert.EqualError(t, err, errormap.ErrItemNotExistOrActive)
}

func TestApproveBorrowing_ItemRepoError(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, mockUserRepo := setupBorrowingServiceMock(t)
	mockDB, sqlMock := setupMockDB(t)

	borrowingId := uint(1)
	approverId := uint(1)
	mockErr := errors.New("database error")

	mockApprover := &entity.User{
		GormModel: entity.GormModel{ID: 1},
		Username:  "test_admin",
		IsAdmin:   true,
	}

	mockBorrowing := &entity.Borrowing{
		GormModel:       entity.GormModel{ID: 1},
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_PENDING,
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	mockItemRepo.EXPECT().GetDB().Return(mockDB)
	mockBorrowingRepo.EXPECT().GetBorrowingByID(borrowingId).Return(mockBorrowing, nil)
	mockUserRepo.EXPECT().GetUserByID(approverId).Return(mockApprover, nil)
	mockItemRepo.EXPECT().GetItemByIDForUpdate(gomock.Any(), mockBorrowing.ItemID).Return(nil, mockErr)

	// act
	result, err := service.ApproveBorrowing(borrowingId, approverId)

	// assert
	assert.Nil(t, result)
	assert.EqualError(t, err, mockErr.Error())
}

func TestApproveBorrowing_NotEnoughItemsForApproval(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, mockUserRepo := setupBorrowingServiceMock(t)
	mockDB, sqlMock := setupMockDB(t)

	borrowingId := uint(1)
	approverId := uint(1)

	mockApprover := &entity.User{
		GormModel: entity.GormModel{ID: 1},
		Username:  "test_admin",
		IsAdmin:   true,
	}

	mockBorrowing := &entity.Borrowing{
		GormModel:       entity.GormModel{ID: 1},
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 5, // Requesting more than available
		BorrowingStatus: entity.BORROWING_PENDING,
	}

	// Item with limited availability
	mockItem := &entity.Item{
		GormModel:       entity.GormModel{ID: 1},
		Name:            "portable projector",
		AvailableAmount: 1,
		TotalAmount:     2,
		Status:          "available",
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	mockItemRepo.EXPECT().GetDB().Return(mockDB)
	mockBorrowingRepo.EXPECT().GetBorrowingByID(borrowingId).Return(mockBorrowing, nil)
	mockUserRepo.EXPECT().GetUserByID(approverId).Return(mockApprover, nil)
	mockItemRepo.EXPECT().GetItemByIDForUpdate(gomock.Any(), mockBorrowing.ItemID).Return(mockItem, nil)

	// act
	result, err := service.ApproveBorrowing(borrowingId, approverId)

	// assert
	assert.Nil(t, result)
	assert.EqualError(t, err, errormap.ErrNotEnoughItemsForApproval)
}

func TestApproveBorrowing_UpdateWithTxError(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, mockUserRepo := setupBorrowingServiceMock(t)
	mockDB, sqlMock := setupMockDB(t)

	borrowingId := uint(1)
	approverId := uint(1)
	mockErr := errors.New("update transaction error")

	mockApprover := &entity.User{
		GormModel: entity.GormModel{ID: 1},
		Username:  "test_admin",
		IsAdmin:   true,
	}

	mockBorrowing := &entity.Borrowing{
		GormModel:       entity.GormModel{ID: 1},
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_PENDING,
	}

	mockItem := &entity.Item{
		GormModel:       entity.GormModel{ID: 1},
		Name:            "office chair ergonomic",
		AvailableAmount: 4,
		TotalAmount:     6,
		Status:          "available",
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	mockItemRepo.EXPECT().GetDB().Return(mockDB)
	mockBorrowingRepo.EXPECT().GetBorrowingByID(borrowingId).Return(mockBorrowing, nil)
	mockUserRepo.EXPECT().GetUserByID(approverId).Return(mockApprover, nil)
	mockItemRepo.EXPECT().GetItemByIDForUpdate(gomock.Any(), mockBorrowing.ItemID).Return(mockItem, nil)
	mockItemRepo.EXPECT().UpdateWithTx(gomock.Any(), gomock.Any()).Return(nil, mockErr)

	// act
	result, err := service.ApproveBorrowing(borrowingId, approverId)

	// assert
	assert.Nil(t, result)
	assert.EqualError(t, err, mockErr.Error())
}

func TestApproveBorrowing_ApproveBorrowingWithTxError(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, mockUserRepo := setupBorrowingServiceMock(t)
	mockDB, sqlMock := setupMockDB(t)

	borrowingId := uint(1)
	approverId := uint(1)
	mockErr := errors.New("approve transaction error")

	mockApprover := &entity.User{
		GormModel: entity.GormModel{ID: 1},
		Username:  "test_admin",
		IsAdmin:   true,
	}

	mockBorrowing := &entity.Borrowing{
		GormModel:       entity.GormModel{ID: 1},
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_PENDING,
	}

	mockItem := &entity.Item{
		GormModel:       entity.GormModel{ID: 1},
		Name:            "office chair ergonomic",
		AvailableAmount: 4,
		TotalAmount:     6,
		Status:          "available",
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	mockItemRepo.EXPECT().GetDB().Return(mockDB)
	mockBorrowingRepo.EXPECT().GetBorrowingByID(borrowingId).Return(mockBorrowing, nil)
	mockUserRepo.EXPECT().GetUserByID(approverId).Return(mockApprover, nil)
	mockItemRepo.EXPECT().GetItemByIDForUpdate(gomock.Any(), mockBorrowing.ItemID).Return(mockItem, nil)
	mockItemRepo.EXPECT().UpdateWithTx(gomock.Any(), gomock.Any()).Return(mockItem, nil)
	mockBorrowingRepo.EXPECT().ApproveBorrowingWithTx(gomock.Any(), borrowingId, approverId).Return(nil, mockErr)

	// act
	result, err := service.ApproveBorrowing(borrowingId, approverId)

	// assert
	assert.Nil(t, result)
	assert.EqualError(t, err, mockErr.Error())
}

func TestRejectBorrowing_Success_AdminReject(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, mockUserRepo := setupBorrowingServiceMock(t)
	mockDB, sqlMock := setupMockDB(t)

	borrowingId := uint(1)
	rejecterId := uint(1)

	mockRejecter := &entity.User{
		GormModel: entity.GormModel{ID: 1},
		Username:  "test_admin",
		IsAdmin:   true,
	}

	mockBorrowing := &entity.Borrowing{
		GormModel:       entity.GormModel{ID: 1},
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_PENDING,
	}

	expectedResult := &entity.Borrowing{
		GormModel:       entity.GormModel{ID: 1},
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_REJECTED,
		ApprovedBy:      rejecterId,
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()

	mockItemRepo.EXPECT().GetDB().Return(mockDB)
	mockBorrowingRepo.EXPECT().GetBorrowingByID(borrowingId).Return(mockBorrowing, nil)
	mockUserRepo.EXPECT().GetUserByID(rejecterId).Return(mockRejecter, nil)
	mockBorrowingRepo.EXPECT().RejectBorrowingWithTx(gomock.Any(), borrowingId, rejecterId).Return(expectedResult, nil)

	// act
	result, err := service.RejectBorrowing(borrowingId, rejecterId)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, entity.BORROWING_REJECTED, result.BorrowingStatus)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestRejectBorrowing_Success_OwnerCancel(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, mockUserRepo := setupBorrowingServiceMock(t)
	mockDB, sqlMock := setupMockDB(t)

	borrowingId := uint(1)
	rejecterId := uint(2) // Same as owner

	mockRejecter := &entity.User{
		GormModel: entity.GormModel{ID: 2},
		Username:  "test_user",
		IsAdmin:   false,
	}

	mockBorrowing := &entity.Borrowing{
		GormModel:       entity.GormModel{ID: 1},
		UserID:          2, // Owner is 2
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_PENDING,
	}

	expectedResult := &entity.Borrowing{
		GormModel:       entity.GormModel{ID: 1},
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_REJECTED,
		ApprovedBy:      rejecterId,
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()

	mockItemRepo.EXPECT().GetDB().Return(mockDB)
	mockBorrowingRepo.EXPECT().GetBorrowingByID(borrowingId).Return(mockBorrowing, nil)
	mockUserRepo.EXPECT().GetUserByID(rejecterId).Return(mockRejecter, nil)
	mockBorrowingRepo.EXPECT().RejectBorrowingWithTx(gomock.Any(), borrowingId, rejecterId).Return(expectedResult, nil)

	// act
	result, err := service.RejectBorrowing(borrowingId, rejecterId)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, entity.BORROWING_REJECTED, result.BorrowingStatus)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestRejectBorrowing_Unauthorized(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, mockUserRepo := setupBorrowingServiceMock(t)
	mockDB, sqlMock := setupMockDB(t)

	borrowingId := uint(1)
	rejecterId := uint(3) // Different user, not admin

	mockRejecter := &entity.User{
		GormModel: entity.GormModel{ID: 3},
		Username:  "other_user",
		IsAdmin:   false,
	}

	mockBorrowing := &entity.Borrowing{
		GormModel:       entity.GormModel{ID: 1},
		UserID:          2, // Owner is 2
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_PENDING,
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	mockItemRepo.EXPECT().GetDB().Return(mockDB)
	mockBorrowingRepo.EXPECT().GetBorrowingByID(borrowingId).Return(mockBorrowing, nil)
	mockUserRepo.EXPECT().GetUserByID(rejecterId).Return(mockRejecter, nil)

	// act
	result, err := service.RejectBorrowing(borrowingId, rejecterId)

	// assert
	assert.Nil(t, result)
	assert.EqualError(t, err, errormap.ErrUnauthorizedToRejectAndCancel)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestRejectBorrowing_BorrowingNotExist(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, _ := setupBorrowingServiceMock(t)
	mockDB, sqlMock := setupMockDB(t)

	borrowingId := uint(999)
	rejecterId := uint(1)

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	mockItemRepo.EXPECT().GetDB().Return(mockDB)
	mockBorrowingRepo.EXPECT().GetBorrowingByID(borrowingId).Return(nil, nil)

	// act
	result, err := service.RejectBorrowing(borrowingId, rejecterId)

	// assert
	assert.Nil(t, result)
	assert.EqualError(t, err, errormap.ErrBorrowingNotExist)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestRejectBorrowing_RejecterNotExist(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, mockUserRepo := setupBorrowingServiceMock(t)
	mockDB, sqlMock := setupMockDB(t)

	borrowingId := uint(1)
	rejecterId := uint(999)

	mockBorrowing := &entity.Borrowing{
		GormModel:       entity.GormModel{ID: 1},
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_PENDING,
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	mockItemRepo.EXPECT().GetDB().Return(mockDB)
	mockBorrowingRepo.EXPECT().GetBorrowingByID(borrowingId).Return(mockBorrowing, nil)
	mockUserRepo.EXPECT().GetUserByID(rejecterId).Return(nil, nil)

	// act
	result, err := service.RejectBorrowing(borrowingId, rejecterId)

	// assert
	assert.Nil(t, result)
	assert.EqualError(t, err, errormap.ErrRejecterNotExist)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestRejectBorrowing_BorrowingNotPending(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, mockUserRepo := setupBorrowingServiceMock(t)
	mockDB, sqlMock := setupMockDB(t)

	borrowingId := uint(1)
	rejecterId := uint(1)

	mockRejecter := &entity.User{
		GormModel: entity.GormModel{ID: 1},
		Username:  "test_admin",
		IsAdmin:   true,
	}

	mockBorrowing := &entity.Borrowing{
		GormModel:       entity.GormModel{ID: 1},
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_ACTIVE,
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	mockItemRepo.EXPECT().GetDB().Return(mockDB)
	mockBorrowingRepo.EXPECT().GetBorrowingByID(borrowingId).Return(mockBorrowing, nil)
	mockUserRepo.EXPECT().GetUserByID(rejecterId).Return(mockRejecter, nil)

	// act
	result, err := service.RejectBorrowing(borrowingId, rejecterId)

	// assert
	assert.Nil(t, result)
	assert.EqualError(t, err, errormap.ErrBorrowingNotPending)
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestRejectBorrowing_RepoError(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, mockUserRepo := setupBorrowingServiceMock(t)
	mockDB, sqlMock := setupMockDB(t)

	borrowingId := uint(1)
	rejecterId := uint(1)
	mockErr := errors.New("database error")

	mockRejecter := &entity.User{
		GormModel: entity.GormModel{ID: 1},
		Username:  "test_admin",
		IsAdmin:   true,
	}

	mockBorrowing := &entity.Borrowing{
		GormModel:       entity.GormModel{ID: 1},
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_PENDING,
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	mockItemRepo.EXPECT().GetDB().Return(mockDB)
	mockBorrowingRepo.EXPECT().GetBorrowingByID(borrowingId).Return(mockBorrowing, nil)
	mockUserRepo.EXPECT().GetUserByID(rejecterId).Return(mockRejecter, nil)
	mockBorrowingRepo.EXPECT().RejectBorrowingWithTx(gomock.Any(), borrowingId, rejecterId).Return(nil, mockErr)

	// act
	result, err := service.RejectBorrowing(borrowingId, rejecterId)

	// assert
	assert.Nil(t, result)
	assert.EqualError(t, err, mockErr.Error())
	assert.NoError(t, sqlMock.ExpectationsWereMet())
}

func TestGetBorrowingsByBorrowingStatus_Success(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, _, _ := setupBorrowingServiceMock(t)

	status := []string{entity.BORROWING_PENDING}
	search := ""
	page := 1
	limit := 10
	expectedBorrowings := &entity.PaginationResult[entity.Borrowing]{
		Data: []entity.Borrowing{
			{
				BorrowingStatus: entity.BORROWING_PENDING,
			},
		},
		TotalItems: 1,
		TotalPages: 1,
		Page:       1,
		Limit:      10,
	}

	mockBorrowingRepo.EXPECT().GetBorrowingsByBorrowingStatus(status, search, page, limit).Return(expectedBorrowings, nil)

	// act
	result, err := service.GetBorrowingsByBorrowingStatus(status, search, page, limit)

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedBorrowings, result)
}

func TestGetBorrowingsByUserID_Success(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, _, _ := setupBorrowingServiceMock(t)

	id := uint(1)
	page := 1
	limit := 10
	search := ""

	expectedBorrowings := &entity.PaginationResult[entity.Borrowing]{
		Data:       []entity.Borrowing{},
		TotalItems: 0,
		TotalPages: 0,
		Page:       1,
		Limit:      10,
	}

	mockBorrowingRepo.EXPECT().GetBorrowingsByUserID(id, page, limit, search).Return(expectedBorrowings, nil)

	// act
	result, err := service.GetBorrowingsByUserID(id, page, limit, search)

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedBorrowings, result)
}

func TestGetBorrowingsStats_Success(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, _, _ := setupBorrowingServiceMock(t)

	id := uint(1)

	expectedBorrowings := &entity.BorrowingStats{}

	mockBorrowingRepo.EXPECT().GetUserBorrowingStats(id).Return(expectedBorrowings, nil)

	// act
	result, err := service.GetUserBorrowingStats(id)

	// assert
	assert.Nil(t, err)
	assert.Equal(t, expectedBorrowings, result)
}

func TestReturnBorrowing_Success(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, _ := setupBorrowingServiceMock(t)
	mockDB, sqlMock := setupMockDB(t)

	borrowingId := uint(1)

	mockBorrowing := &entity.Borrowing{
		GormModel:       entity.GormModel{ID: 1},
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_ACTIVE,
		DueDate:         time.Now().Add(24 * time.Hour).Format(time.RFC3339),
	}

	mockItem := &entity.Item{
		GormModel:       entity.GormModel{ID: 1},
		Name:            "office chair ergonomic",
		AvailableAmount: 4,
		TotalAmount:     6,
	}

	expectedResult := &entity.Borrowing{
		GormModel:       entity.GormModel{ID: 1},
		BorrowingStatus: entity.BORROWING_RETURNED,
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()

	mockItemRepo.EXPECT().GetDB().Return(mockDB)
	mockBorrowingRepo.EXPECT().GetBorrowingByID(borrowingId).Return(mockBorrowing, nil)
	mockItemRepo.EXPECT().GetItemByIDForUpdate(gomock.Any(), mockBorrowing.ItemID).Return(mockItem, nil)
	mockItemRepo.EXPECT().UpdateWithTx(gomock.Any(), gomock.Any()).Return(mockItem, nil)
	mockBorrowingRepo.EXPECT().ReturnBorrowingWithTx(gomock.Any(), borrowingId, entity.BORROWING_RETURNED).Return(expectedResult, nil)

	// act
	result, err := service.ReturnBorrowing(borrowingId)

	// assert
	assert.Nil(t, err)
	assert.Equal(t, entity.BORROWING_RETURNED, result.BorrowingStatus)
}

func TestReturnBorrowing_Success_Overdue(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, _ := setupBorrowingServiceMock(t)
	mockDB, sqlMock := setupMockDB(t)

	borrowingId := uint(1)

	// Due date in the past
	mockBorrowing := &entity.Borrowing{
		GormModel:       entity.GormModel{ID: 1},
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_ACTIVE,
		DueDate:         time.Now().Add(-24 * time.Hour).Format(time.RFC3339),
	}

	mockItem := &entity.Item{
		GormModel:       entity.GormModel{ID: 1},
		Name:            "office chair ergonomic",
		AvailableAmount: 4,
	}

	expectedResult := &entity.Borrowing{
		GormModel:       entity.GormModel{ID: 1},
		BorrowingStatus: entity.BORROWING_OVERDUE,
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectCommit()

	mockItemRepo.EXPECT().GetDB().Return(mockDB)
	mockBorrowingRepo.EXPECT().GetBorrowingByID(borrowingId).Return(mockBorrowing, nil)
	mockItemRepo.EXPECT().GetItemByIDForUpdate(gomock.Any(), mockBorrowing.ItemID).Return(mockItem, nil)
	mockItemRepo.EXPECT().UpdateWithTx(gomock.Any(), gomock.Any()).Return(mockItem, nil)
	mockBorrowingRepo.EXPECT().ReturnBorrowingWithTx(gomock.Any(), borrowingId, entity.BORROWING_OVERDUE).Return(expectedResult, nil)

	// act
	result, err := service.ReturnBorrowing(borrowingId)

	// assert
	assert.Nil(t, err)
	assert.Equal(t, entity.BORROWING_OVERDUE, result.BorrowingStatus)
}

func TestReturnBorrowing_NotExist(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, _ := setupBorrowingServiceMock(t)
	mockDB, sqlMock := setupMockDB(t)

	borrowingId := uint(999)

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	mockItemRepo.EXPECT().GetDB().Return(mockDB)
	mockBorrowingRepo.EXPECT().GetBorrowingByID(borrowingId).Return(nil, nil)

	// act
	result, err := service.ReturnBorrowing(borrowingId)

	// assert
	assert.Nil(t, result)
	assert.EqualError(t, err, errormap.ErrBorrowingNotExist)
}

func TestReturnBorrowing_NotActive(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, _ := setupBorrowingServiceMock(t)
	mockDB, sqlMock := setupMockDB(t)

	borrowingId := uint(1)

	mockBorrowing := &entity.Borrowing{
		GormModel:       entity.GormModel{ID: 1},
		BorrowingStatus: entity.BORROWING_PENDING, // Not Active
	}

	sqlMock.ExpectBegin()
	sqlMock.ExpectRollback()

	mockItemRepo.EXPECT().GetDB().Return(mockDB)
	mockBorrowingRepo.EXPECT().GetBorrowingByID(borrowingId).Return(mockBorrowing, nil)

	// act
	result, err := service.ReturnBorrowing(borrowingId)

	// assert
	assert.Nil(t, result)
	assert.EqualError(t, err, errormap.ErrBorrowingNotActive)
}
