package port

import (
	"github.com/nitikhon/golang-inventory-system/internal/core/entity"
	"gorm.io/gorm"
)

type BorrowingRepository interface {
	BorrowItem(borrowing *entity.Borrowing) (*entity.Borrowing, error)
	ApproveBorrowingWithTx(tx *gorm.DB, borrowingId, approverId uint) (*entity.Borrowing, error)
	RejectBorrowingWithTx(tx *gorm.DB, borrowingId, rejecterId uint) (*entity.Borrowing, error)
	GetAllBorrowings() ([]*entity.Borrowing, error)
	GetBorrowingByID(borrowingID uint) (*entity.Borrowing, error)
	GetBorrowingsByUserID(borrowerId uint, page, limit int, search string) (*entity.PaginationResult[entity.Borrowing], error)
	GetBorrowingsByItemID(itemID uint) ([]*entity.Borrowing, error)
	GetBorrowingsByBorrowingStatus(status []string, search string, page, limit int) (*entity.PaginationResult[entity.Borrowing], error)
	GetBorrowingsByApprovalStatus(status string) ([]*entity.Borrowing, error)
	GetBorrowingsByApproverID(approverID uint) ([]*entity.Borrowing, error)
	GetUserBorrowingStats(userID uint) (*entity.BorrowingStats, error)
	HasActiveBorrowing(userID, itemID uint) (bool, error)
}
