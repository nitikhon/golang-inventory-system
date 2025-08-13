package entity

import (
	"gorm.io/gorm"
)

const (
	// BorrowingStatus constants
	BORROWING_PENDING   = "pending"
	BORROWING_ACTIVE   = "active"
	BORROWING_RETURNED = "returned"
	BORROWING_OVERDUE  = "overdue"
	BORROWING_CANCELLED = "cancelled"
    BORROWING_LOST     = "lost"

	// ApprovalStatus constants
	APPROVAL_PENDING   = "pending"
	APPROVAL_APPROVED  = "approved"
	APPROVAL_REJECTED  = "rejected"
)

type Borrowing struct {
	UserID          uint `gorm:"primaryKey"`
	ItemID          uint `gorm:"primaryKey"`
	Description     string
	BorrowedAt      string
	ReturnedAt      string
	DueDate         string
	BorrowingAmount int
	BorrowingStatus string `gorm:"type:VARCHAR(20);default:'pending';index"`
	ApprovalStatus  string `gorm:"type:VARCHAR(20);default:'pending';index"` // pending, approved, rejected
	ApprovedAt      string
	ApprovedBy      uint
	RejectedBy      uint
	gorm.Model
}
