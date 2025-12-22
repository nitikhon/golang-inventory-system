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
	GetBorrowingsByUserID(userID uint) ([]*entity.Borrowing, error)
	GetBorrowingsByItemID(itemID uint) ([]*entity.Borrowing, error)
	GetBorrowingsByBorrowingStatus(status string) ([]*entity.Borrowing, error)
	GetBorrowingsByApprovalStatus(status string) ([]*entity.Borrowing, error)
	GetBorrowingsByApproverID(approverID uint) ([]*entity.Borrowing, error)
}
