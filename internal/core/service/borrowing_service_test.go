package service

import (
	"errors"
	"testing"

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
		Model:     gorm.Model{ID: 2},
		Username:  "test_user",
		Email:     "user@test.com",
		Phone:     "0987654321",
		FirstName: "Test",
		LastName:  "User",
		IsAdmin:   false,
	}

	mockItem := &entity.Item{
		Model:           gorm.Model{ID: 1},
		Name:            "office chair ergonomic",
		Description:     "Ergonomic office chair with lumbar support and adjustable height",
		AvailableAmount: 4,
		TotalAmount:     6,
		Status:          "available",
	}

	borrowingInput := entity.Borrowing{
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
	}

	expectedBorrowing := &entity.Borrowing{
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_PENDING,
		ApprovalStatus:  entity.APPROVAL_PENDING,
	}

	mockUserRepo.EXPECT().GetUserByID(borrowingInput.UserID).Return(mockUser, nil)
	mockItemRepo.EXPECT().GetItemByID(borrowingInput.ItemID).Return(mockItem, nil)
	mockBorrowingRepo.EXPECT().GetBorrowingsByUserID(borrowingInput.UserID).Return([]*entity.Borrowing{}, nil)
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

	borrowingInput := entity.Borrowing{
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

	borrowingInput := entity.Borrowing{
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
		Model:    gorm.Model{ID: 2},
		Username: "test_user",
	}

	// Item with status "borrowed" (based on setup.go - conference camera 4k)
	mockItem := &entity.Item{
		Model:           gorm.Model{ID: 4},
		Name:            "conference camera 4k",
		AvailableAmount: 0,
		TotalAmount:     4,
		Status:          "borrowed",
	}

	borrowingInput := entity.Borrowing{
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
		Model:    gorm.Model{ID: 2},
		Username: "test_user",
	}

	borrowingInput := entity.Borrowing{
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
		Model:    gorm.Model{ID: 2},
		Username: "test_user",
	}

	borrowingInput := entity.Borrowing{
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
		Model:    gorm.Model{ID: 2},
		Username: "test_user",
	}

	// Item with limited availability (based on setup.go - portable projector)
	mockItem := &entity.Item{
		Model:           gorm.Model{ID: 6},
		Name:            "portable projector",
		AvailableAmount: 1,
		TotalAmount:     2,
		Status:          "available",
	}

	borrowingInput := entity.Borrowing{
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
		Model:    gorm.Model{ID: 2},
		Username: "test_user",
	}

	mockItem := &entity.Item{
		Model:           gorm.Model{ID: 1},
		Name:            "office chair ergonomic",
		AvailableAmount: 4,
		TotalAmount:     6,
		Status:          "available",
	}

	// Existing pending borrowing for the same item
	existingBorrowing := []*entity.Borrowing{
		{
			UserID:          2,
			ItemID:          1,
			BorrowingStatus: entity.BORROWING_PENDING,
		},
	}

	borrowingInput := entity.Borrowing{
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
	}

	mockUserRepo.EXPECT().GetUserByID(borrowingInput.UserID).Return(mockUser, nil)
	mockItemRepo.EXPECT().GetItemByID(borrowingInput.ItemID).Return(mockItem, nil)
	mockBorrowingRepo.EXPECT().GetBorrowingsByUserID(borrowingInput.UserID).Return(existingBorrowing, nil)

	// act
	result, err := service.BorrowItem(borrowingInput)

	// assert
	assert.Nil(t, result)
	assert.EqualError(t, err, errormap.ErrAlreadyBorrowed)
}

func TestBorrowItem_GetBorrowingsRepoError(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, mockUserRepo := setupBorrowingServiceMock(t)

	mockUser := &entity.User{
		Model:    gorm.Model{ID: 2},
		Username: "test_user",
	}

	mockItem := &entity.Item{
		Model:           gorm.Model{ID: 1},
		Name:            "office chair ergonomic",
		AvailableAmount: 4,
		TotalAmount:     6,
		Status:          "available",
	}

	borrowingInput := entity.Borrowing{
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
	}
	mockErr := errors.New("database error")

	mockUserRepo.EXPECT().GetUserByID(borrowingInput.UserID).Return(mockUser, nil)
	mockItemRepo.EXPECT().GetItemByID(borrowingInput.ItemID).Return(mockItem, nil)
	mockBorrowingRepo.EXPECT().GetBorrowingsByUserID(borrowingInput.UserID).Return(nil, mockErr)

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
		Model:    gorm.Model{ID: 2},
		Username: "test_user",
	}

	mockItem := &entity.Item{
		Model:           gorm.Model{ID: 1},
		Name:            "office chair ergonomic",
		AvailableAmount: 4,
		TotalAmount:     6,
		Status:          "available",
	}

	borrowingInput := entity.Borrowing{
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
	}
	mockErr := errors.New("database error")

	mockUserRepo.EXPECT().GetUserByID(borrowingInput.UserID).Return(mockUser, nil)
	mockItemRepo.EXPECT().GetItemByID(borrowingInput.ItemID).Return(mockItem, nil)
	mockBorrowingRepo.EXPECT().GetBorrowingsByUserID(borrowingInput.UserID).Return([]*entity.Borrowing{}, nil)
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
		Model:    gorm.Model{ID: 2},
		Username: "test_user",
	}

	// First item - wireless keyboard (based on setup.go)
	mockItem := &entity.Item{
		Model:           gorm.Model{ID: 3},
		Name:            "wireless keyboard",
		AvailableAmount: 7,
		TotalAmount:     10,
		Status:          "available",
	}

	// User has pending borrowing for a DIFFERENT item (item 1)
	existingBorrowing := []*entity.Borrowing{
		{
			UserID:          2,
			ItemID:          1, // Different item
			BorrowingStatus: entity.BORROWING_PENDING,
		},
	}

	borrowingInput := entity.Borrowing{
		UserID:          2,
		ItemID:          3, // Different item
		BorrowingAmount: 2,
	}

	expectedBorrowing := &entity.Borrowing{
		UserID:          2,
		ItemID:          3,
		BorrowingAmount: 2,
		BorrowingStatus: entity.BORROWING_PENDING,
		ApprovalStatus:  entity.APPROVAL_PENDING,
	}

	mockUserRepo.EXPECT().GetUserByID(borrowingInput.UserID).Return(mockUser, nil)
	mockItemRepo.EXPECT().GetItemByID(borrowingInput.ItemID).Return(mockItem, nil)
	mockBorrowingRepo.EXPECT().GetBorrowingsByUserID(borrowingInput.UserID).Return(existingBorrowing, nil)
	mockBorrowingRepo.EXPECT().BorrowItem(gomock.Any()).Return(expectedBorrowing, nil)

	// act
	result, err := service.BorrowItem(borrowingInput)

	// assert
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, expectedBorrowing.ItemID, result.ItemID)
}

func TestBorrowItem_CanBorrowIfPreviousBorrowingNotPending(t *testing.T) {
	// arrange
	service, mockBorrowingRepo, mockItemRepo, mockUserRepo := setupBorrowingServiceMock(t)

	mockUser := &entity.User{
		Model:    gorm.Model{ID: 2},
		Username: "test_user",
	}

	mockItem := &entity.Item{
		Model:           gorm.Model{ID: 1},
		Name:            "office chair ergonomic",
		AvailableAmount: 4,
		TotalAmount:     6,
		Status:          "available",
	}

	// User has a returned borrowing for the same item (not pending)
	existingBorrowing := []*entity.Borrowing{
		{
			UserID:          2,
			ItemID:          1,
			BorrowingStatus: entity.BORROWING_RETURNED, // Already returned
		},
	}

	borrowingInput := entity.Borrowing{
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
	}

	expectedBorrowing := &entity.Borrowing{
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_PENDING,
		ApprovalStatus:  entity.APPROVAL_PENDING,
	}

	mockUserRepo.EXPECT().GetUserByID(borrowingInput.UserID).Return(mockUser, nil)
	mockItemRepo.EXPECT().GetItemByID(borrowingInput.ItemID).Return(mockItem, nil)
	mockBorrowingRepo.EXPECT().GetBorrowingsByUserID(borrowingInput.UserID).Return(existingBorrowing, nil)
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
		Model:    gorm.Model{ID: 1},
		Username: "test_admin",
		IsAdmin:  true,
	}

	mockBorrowing := &entity.Borrowing{
		Model:           gorm.Model{ID: 1},
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_PENDING,
		ApprovalStatus:  entity.APPROVAL_PENDING,
	}

	mockItem := &entity.Item{
		Model:           gorm.Model{ID: 1},
		Name:            "office chair ergonomic",
		AvailableAmount: 4,
		TotalAmount:     6,
		Status:          "available",
	}

	expectedResult := &entity.Borrowing{
		Model:           gorm.Model{ID: 1},
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_ACTIVE,
		ApprovalStatus:  entity.APPROVAL_APPROVED,
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
	assert.Equal(t, entity.APPROVAL_APPROVED, result.ApprovalStatus)
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
		Model:           gorm.Model{ID: 1},
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_PENDING,
		ApprovalStatus:  entity.APPROVAL_PENDING,
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
		Model:           gorm.Model{ID: 1},
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_PENDING,
		ApprovalStatus:  entity.APPROVAL_PENDING,
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
		Model:    gorm.Model{ID: 1},
		Username: "test_admin",
		IsAdmin:  true,
	}

	// Borrowing already approved
	mockBorrowing := &entity.Borrowing{
		Model:           gorm.Model{ID: 1},
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_ACTIVE,
		ApprovalStatus:  entity.APPROVAL_APPROVED,
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
		Model:    gorm.Model{ID: 1},
		Username: "test_admin",
		IsAdmin:  true,
	}

	// Borrowing already rejected
	mockBorrowing := &entity.Borrowing{
		Model:           gorm.Model{ID: 1},
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_PENDING,
		ApprovalStatus:  entity.APPROVAL_REJECTED,
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
		Model:    gorm.Model{ID: 1},
		Username: "test_admin",
		IsAdmin:  true,
	}

	mockBorrowing := &entity.Borrowing{
		Model:           gorm.Model{ID: 1},
		UserID:          2,
		ItemID:          999, // Non-existent item
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_PENDING,
		ApprovalStatus:  entity.APPROVAL_PENDING,
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
		Model:    gorm.Model{ID: 1},
		Username: "test_admin",
		IsAdmin:  true,
	}

	mockBorrowing := &entity.Borrowing{
		Model:           gorm.Model{ID: 1},
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_PENDING,
		ApprovalStatus:  entity.APPROVAL_PENDING,
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
		Model:    gorm.Model{ID: 1},
		Username: "test_admin",
		IsAdmin:  true,
	}

	mockBorrowing := &entity.Borrowing{
		Model:           gorm.Model{ID: 1},
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 5, // Requesting more than available
		BorrowingStatus: entity.BORROWING_PENDING,
		ApprovalStatus:  entity.APPROVAL_PENDING,
	}

	// Item with limited availability
	mockItem := &entity.Item{
		Model:           gorm.Model{ID: 1},
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
		Model:    gorm.Model{ID: 1},
		Username: "test_admin",
		IsAdmin:  true,
	}

	mockBorrowing := &entity.Borrowing{
		Model:           gorm.Model{ID: 1},
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_PENDING,
		ApprovalStatus:  entity.APPROVAL_PENDING,
	}

	mockItem := &entity.Item{
		Model:           gorm.Model{ID: 1},
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
		Model:    gorm.Model{ID: 1},
		Username: "test_admin",
		IsAdmin:  true,
	}

	mockBorrowing := &entity.Borrowing{
		Model:           gorm.Model{ID: 1},
		UserID:          2,
		ItemID:          1,
		BorrowingAmount: 1,
		BorrowingStatus: entity.BORROWING_PENDING,
		ApprovalStatus:  entity.APPROVAL_PENDING,
	}

	mockItem := &entity.Item{
		Model:           gorm.Model{ID: 1},
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
