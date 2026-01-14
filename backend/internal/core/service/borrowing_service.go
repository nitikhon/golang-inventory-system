package service

import (
	"context"
	"errors"
	"time"

	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"github.com/nitikhon/golang-inventory-system/internal/core/port"
	"github.com/nitikhon/golang-inventory-system/internal/util/errormap"
)

// Ensures BorrowingService imnplements BorrowingServiceInterface
var _ BorrowingServiceInterface = (*BorrowingService)(nil)

// BorrowingServiceInterface defines the contract for borrowing operations.
type BorrowingServiceInterface interface {
	BorrowItem(ctx context.Context, borrowing entity.BorrowRequest) (*entity.Borrowing, error)
	ApproveBorrowing(ctx context.Context, borrowerId, approverId uint) (*entity.Borrowing, error)
	RejectBorrowing(ctx context.Context, borrowerId, rejecterId uint) (*entity.Borrowing, error)
	ReturnBorrowing(ctx context.Context, borrowerId uint) (*entity.Borrowing, error)
	GetBorrowingsByBorrowingStatus(ctx context.Context, status []string, search string, page, limit int) (*entity.PaginationResult[entity.Borrowing], error)
	GetBorrowingsByUserID(ctx context.Context, borrowerId uint, page, limit int, search string) (*entity.PaginationResult[entity.Borrowing], error)
	GetUserBorrowingStats(ctx context.Context, userID uint) (*entity.BorrowingStats, error)
	MarkOverdueItems(ctx context.Context) error
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
func (s *BorrowingService) BorrowItem(ctx context.Context, req entity.BorrowRequest) (*entity.Borrowing, error) {
	// User exists and is active
	user, err := s.userRepo.GetUserByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if user == nil || !user.DeletedAt.Time.IsZero() {
		return nil, errors.New(errormap.ErrUserNotExist)
	}

	// Item exists and is available
	item, err := s.itemRepo.GetItemByID(ctx, req.ItemID)
	if err != nil {
		return nil, err
	}
	if item == nil || !item.DeletedAt.Time.IsZero() || item.Status != "available" {
		return nil, errors.New(errormap.ErrItemNotAvailable)
	}

	// Check available quantity
	if item.AvailableAmount < req.BorrowingAmount {
		return nil, errors.New(errormap.ErrItemNotEnough)
	}

	// Check if the user has already borrowed this item
	exists, err := s.borrowingRepo.HasActiveBorrowing(ctx, req.UserID, req.ItemID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New(errormap.ErrAlreadyBorrowed)
	}

	borrowsAt := time.Now()

	if req.DueDate == "" {
		req.DueDate = time.Now().AddDate(0, 0, 7).Format(time.RFC3339) // Default to 7 days later
	}

	if req.DueDate != "" {
		due, err1 := time.Parse(time.RFC3339, req.DueDate)

		if err1 != nil {
			return nil, errors.New(errormap.ErrInvalidDueDateFormat)
		}

		if due.Before(borrowsAt) {
			return nil, errors.New(errormap.ErrInvalidDueDateValue)
		}
	}

	borrowing := entity.Borrowing{
		UserID:          req.UserID,
		ItemID:          req.ItemID,
		Description:     req.Description,
		BorrowingAmount: req.BorrowingAmount,
		BorrowedAt:      borrowsAt.Format(time.RFC3339),
		DueDate:         req.DueDate,
		BorrowingStatus: "pending",
	}

	return s.borrowingRepo.BorrowItem(ctx, &borrowing)
}

// ApproveBorrowing approves a borrowing request.
func (s *BorrowingService) ApproveBorrowing(ctx context.Context, borrowerId, approverId uint) (*entity.Borrowing, error) {
	db := s.itemRepo.GetDB()
	tx := db.WithContext(ctx).Begin()
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
	existingBorrowing, err := s.borrowingRepo.GetBorrowingByID(ctx, borrowerId)
	if err != nil {
		return nil, err
	}
	if existingBorrowing == nil {
		return nil, errors.New(errormap.ErrBorrowingNotExist)
	}

	// Check if the approver exists
	approver, err := s.userRepo.GetUserByID(ctx, approverId)
	if err != nil {
		return nil, err
	}
	if approver == nil || !approver.DeletedAt.Time.IsZero() {
		return nil, errors.New(errormap.ErrApproverNotExist)
	}

	// Check if the borrowing is already approved or rejected or the borrowing status is not pending
	if existingBorrowing.BorrowingStatus != entity.BORROWING_PENDING {
		return nil, errors.New(errormap.ErrBorrowingNotPending)
	}

	// Decrease the item's available amount
	item, err := s.itemRepo.GetItemByIDForUpdate(tx, existingBorrowing.ItemID)
	if err != nil {
		return nil, err
	}
	if item == nil || !item.DeletedAt.Time.IsZero() {
		return nil, errors.New(errormap.ErrItemNotExistOrActive)
	}
	if item.AvailableAmount < existingBorrowing.BorrowingAmount {
		return nil, errors.New(errormap.ErrNotEnoughItemsForApproval)
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
func (s *BorrowingService) RejectBorrowing(ctx context.Context, borrowerId, rejecterId uint) (*entity.Borrowing, error) {
	db := s.itemRepo.GetDB()
	tx := db.WithContext(ctx).Begin()
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
	existingBorrowing, err := s.borrowingRepo.GetBorrowingByID(ctx, borrowerId)
	if err != nil {
		return nil, err
	}
	if existingBorrowing == nil {
		return nil, errors.New(errormap.ErrBorrowingNotExist)
	}

	// Check if the rejecter exists
	rejecter, err := s.userRepo.GetUserByID(ctx, rejecterId)
	if err != nil {
		return nil, err
	}
	if rejecter == nil || !rejecter.DeletedAt.Time.IsZero() {
		return nil, errors.New(errormap.ErrRejecterNotExist)
	}

	// Check if the rejecter is owner of the borrowing or admin
	if rejecter.ID != existingBorrowing.UserID && !rejecter.IsAdmin {
		return nil, errors.New(errormap.ErrUnauthorizedToRejectAndCancel)
	}

	// Check if the borrowing is already approved or rejected or the borrowing status is not pending
	if existingBorrowing.BorrowingStatus != entity.BORROWING_PENDING {
		return nil, errors.New(errormap.ErrBorrowingNotPending)
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

// RejectBorrowing rejects a borrowing request.
func (s *BorrowingService) ReturnBorrowing(ctx context.Context, borrowingId uint) (*entity.Borrowing, error) {
	db := s.itemRepo.GetDB()
	tx := db.WithContext(ctx).Begin()
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
	existingBorrowing, err := s.borrowingRepo.GetBorrowingByID(ctx, borrowingId)
	if err != nil {
		return nil, err
	}
	if existingBorrowing == nil {
		return nil, errors.New(errormap.ErrBorrowingNotExist)
	}

	// Check if the borrowing is active or not
	if existingBorrowing.BorrowingStatus != entity.BORROWING_ACTIVE {
		return nil, errors.New(errormap.ErrBorrowingNotActive)
	}

	now := time.Now()
	status := entity.BORROWING_RETURNED

	// assume that we don't have legacy data with wrong format
	dueDate, _ := time.Parse(time.RFC3339, existingBorrowing.DueDate)

	currentDate := now.Truncate(24 * time.Hour)
	dueDateDate := dueDate.Truncate(24 * time.Hour)
	if currentDate.After(dueDateDate) {
		status = entity.BORROWING_OVERDUE
	}

	item, err := s.itemRepo.GetItemByIDForUpdate(tx, existingBorrowing.ItemID)
	if err != nil {
		return nil, err
	}
	item.AvailableAmount += existingBorrowing.BorrowingAmount

	if _, err := s.itemRepo.UpdateWithTx(tx, item); err != nil {
		tx.Rollback()
		return nil, err
	}

	result, err := s.borrowingRepo.ReturnBorrowingWithTx(tx, borrowingId, status)
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
func (s *BorrowingService) GetBorrowingsByBorrowingStatus(ctx context.Context, status []string, search string, page, limit int) (*entity.PaginationResult[entity.Borrowing], error) {
	return s.borrowingRepo.GetBorrowingsByBorrowingStatus(ctx, status, search, page, limit)
}

func (s *BorrowingService) GetBorrowingsByUserID(ctx context.Context, id uint, page, limit int, search string) (*entity.PaginationResult[entity.Borrowing], error) {
	return s.borrowingRepo.GetBorrowingsByUserID(ctx, id, page, limit, search)
}

func (s *BorrowingService) GetUserBorrowingStats(ctx context.Context, userID uint) (*entity.BorrowingStats, error) {
	return s.borrowingRepo.GetUserBorrowingStats(ctx, userID)
}

func (s *BorrowingService) MarkOverdueItems(ctx context.Context) error {
	db := s.borrowingRepo.GetDB()
	tx := db.WithContext(ctx).Begin()
	committed := false
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
		if !committed {
			tx.Rollback()
		}
	}()

	err := s.borrowingRepo.MarkOverdueItemsWithTx(tx)

	if err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		return err
	}
	committed = true

	return nil
}
