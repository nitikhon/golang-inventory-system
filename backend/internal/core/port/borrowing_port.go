package port

import (
	"context"

	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"gorm.io/gorm"
)

type BorrowingRepository interface {
	BorrowItem(ctx context.Context, borrowing *entity.Borrowing) (*entity.Borrowing, error)
	ApproveBorrowingWithTx(tx *gorm.DB, borrowingId, approverId uint) (*entity.Borrowing, error)
	RejectBorrowingWithTx(tx *gorm.DB, borrowingId, rejecterId uint) (*entity.Borrowing, error)
	ReturnBorrowingWithTx(tx *gorm.DB, borrowingId uint, status string) (*entity.Borrowing, error)
	GetAllBorrowings(ctx context.Context) ([]*entity.Borrowing, error)
	GetBorrowingByID(ctx context.Context, borrowingID uint) (*entity.Borrowing, error)
	GetBorrowingsByUserID(ctx context.Context, borrowerId uint, page, limit int, search string) (*entity.PaginationResult[entity.Borrowing], error)
	GetBorrowingsByItemID(ctx context.Context, itemID uint) ([]*entity.Borrowing, error)
	GetBorrowingsByBorrowingStatus(ctx context.Context, status []string, search string, page, limit int) (*entity.PaginationResult[entity.Borrowing], error)
	GetBorrowingsByApproverID(ctx context.Context, approverID uint) ([]*entity.Borrowing, error)
	GetUserBorrowingStats(ctx context.Context, userID uint) (*entity.BorrowingStats, error)
	HasActiveBorrowing(ctx context.Context, userID, itemID uint) (bool, error)
	MarkOverdueItemsWithTx(tx *gorm.DB) error
	GetDB() *gorm.DB
}
