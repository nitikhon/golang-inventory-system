package service

import (
	"errors"
	"time"

	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"github.com/nitikhon/golang-inventory-system/internal/core/port"
)

// Ensures BorrowingService imnplements BorrowingServiceInterface
var _ BorrowingServiceInterface = (*BorrowingService)(nil)

// BorrowingServiceInterface defines the contract for borrowing operations.
type BorrowingServiceInterface interface {
	BorrowItem(borrowing entity.Borrowing) (*entity.Borrowing, error)
	ApproveBorrowing(borrowerId, approverId uint) (*entity.Borrowing, error)
	RejectBorrowing(borrowerId, rejecterId uint) (*entity.Borrowing, error)
	GetBorrowingsByBorrowingStatus(status string) ([]*entity.Borrowing, error)
	GetBorrowingsByApprovalStatus(status string) ([]*entity.Borrowing, error)
}

// BorrowingService provides the use cases for borrowing operations.
type BorrowingService struct {
	borrowingRepo port.BorrowingRepository
	itemRepo      port.ItemRepository
	userRepo      port.UserRepository
}

// NewBorrowingService creates a new BorrowingService instance.
func NewBorrowingService(borrowingRepo port.BorrowingRepository, itemRepo port.ItemRepository, userRepo port.UserRepository) *BorrowingService {
	return &BorrowingService{
		borrowingRepo: borrowingRepo,
		itemRepo:      itemRepo,
		userRepo:      userRepo,
	}
}

// BorrowItem handles the borrowing of an item.
func (s *BorrowingService) BorrowItem(borrowing entity.Borrowing) (*entity.Borrowing, error) {
	// 1. User exists and is active
	user, err := s.userRepo.GetUserByID(borrowing.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil || !user.DeletedAt.Time.IsZero() {
		return nil, errors.New("user does not exist or is not active")
	}

	// 2. Item exists and is available
	item, err := s.itemRepo.GetItemByID(borrowing.ItemID)
	if err != nil {
		return nil, err
	}
	if item == nil || !item.DeletedAt.Time.IsZero() || item.Status != "available" {
		return nil, errors.New("item is not available for borrowing")
	}

	// 3. Check available quantity
	if item.AvailableAmount < borrowing.BorrowingAmount {
		return nil, errors.New("item is not available enough to borrow")
	}

	// 4. Check if the user has already borrowed this item
	existingBorrowings, err := s.borrowingRepo.GetBorrowingsByUserID(borrowing.UserID)
	if err != nil {
		return nil, err
	}
	for _, existing := range existingBorrowings {
		if existing.ItemID == borrowing.ItemID && existing.BorrowingStatus == "pending" {
			return nil, errors.New("user has already borrowed this item and it is still pending")
		}
	}

	// 5. Set default dates if not provided
	if borrowing.BorrowedAt == "" {
		borrowing.BorrowedAt = time.Now().Format(time.RFC3339)
	}

	if borrowing.DueDate == "" {
		borrowing.DueDate = time.Now().AddDate(0, 0, 7).Format(time.RFC3339) // Default to 7 days later
	}

	return s.borrowingRepo.BorrowItem(&borrowing)
}

// ApproveBorrowing approves a borrowing request.
func (s *BorrowingService) ApproveBorrowing(borrowerId, approverId uint) (*entity.Borrowing, error) {
	db := s.itemRepo.GetDB()
	tx := db.Begin()
	committed := false
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
		if !committed {
			tx.Rollback()
		}
	}()

	// Check if the borrowing exists
	existingBorrowing, err := s.borrowingRepo.GetBorrowingByID(borrowerId)
	if err != nil {
		return nil, err
	}
	if existingBorrowing == nil {
		return nil, errors.New("borrowing does not exist")
	}

	// Check if the approver exists
	approver, err := s.userRepo.GetUserByID(approverId)
	if err != nil {
		return nil, err
	}
	if approver == nil || !approver.DeletedAt.Time.IsZero() {
		return nil, errors.New("approver does not exist or is not active")
	}

	// Check if the borrowing is already approved or rejected or the borrowing status is not pending
	if existingBorrowing.ApprovalStatus != entity.APPROVAL_PENDING ||
		existingBorrowing.BorrowingStatus != entity.BORROWING_PENDING {
		return nil, errors.New("borrowing request is not pending")
	}

	// Decrease the item's available amount
	item, err := s.itemRepo.GetItemByIDForUpdate(tx, existingBorrowing.ItemID)
	if err != nil {
		return nil, err
	}
	if item == nil || !item.DeletedAt.Time.IsZero() {
		return nil, errors.New("item does not exist or is not active")
	}
	if item.AvailableAmount < existingBorrowing.BorrowingAmount {
		return nil, errors.New("not enough items available to approve borrowing")
	}
	item.AvailableAmount -= existingBorrowing.BorrowingAmount

	if _, err := s.itemRepo.UpdateWithTx(tx, item); err != nil {
		tx.Rollback()
		return nil, err
	}
	result, err := s.borrowingRepo.ApproveBorrowingWithTx(tx, borrowerId, approverId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	committed = true
	return result, nil
}

// RejectBorrowing rejects a borrowing request.
func (s *BorrowingService) RejectBorrowing(borrowerId, rejecterId uint) (*entity.Borrowing, error) {
	db := s.itemRepo.GetDB()
	tx := db.Begin()
	committed := false
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
		if !committed {
			tx.Rollback()
		}
	}()

	// Check if the borrowing exists
	existingBorrowing, err := s.borrowingRepo.GetBorrowingByID(borrowerId)
	if err != nil {
		return nil, err
	}
	if existingBorrowing == nil {
		return nil, errors.New("borrowing does not exist")
	}

	// Check if the rejecter exists
	rejecter, err := s.userRepo.GetUserByID(rejecterId)
	if err != nil {
		return nil, err
	}
	if rejecter == nil || !rejecter.DeletedAt.Time.IsZero() {
		return nil, errors.New("rejecter does not exist or is not active")
	}

	// Check if the borrowing is already approved or rejected or the borrowing status is not pending
	if existingBorrowing.ApprovalStatus != entity.APPROVAL_PENDING ||
		existingBorrowing.BorrowingStatus != entity.BORROWING_PENDING {
		return nil, errors.New("borrowing request is not pending")
	}

	result, err := s.borrowingRepo.RejectBorrowingWithTx(tx, borrowerId, rejecterId)
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	if err := tx.Commit().Error; err != nil {
		return nil, err
	}
	committed = true
	return result, nil
}

// GetBorrowingByStatus retrieves borrowings by their status.
func (s *BorrowingService) GetBorrowingsByBorrowingStatus(status string) ([]*entity.Borrowing, error) {
	return s.borrowingRepo.GetBorrowingsByBorrowingStatus(status)
}

// GetBorrowingsByApprovalStatus retrieves borrowings by their approval status.
func (s *BorrowingService) GetBorrowingsByApprovalStatus(status string) ([]*entity.Borrowing, error) {
	return s.borrowingRepo.GetBorrowingsByApprovalStatus(status)
}
