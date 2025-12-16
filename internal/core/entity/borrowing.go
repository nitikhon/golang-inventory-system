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
	UserID          uint `json:"user_id" gorm:"primaryKey"`
	ItemID          uint `json:"item_id" gorm:"primaryKey"`
	Description     string `json:"description"`
	BorrowedAt      string `json:"borrowed_at"`
	ReturnedAt      string `json:"returned_at"`
	DueDate         string `json:"due_date"`
	BorrowingAmount int    `json:"borrowing_amount"`
	BorrowingStatus string `json:"borrowing_status" gorm:"type:VARCHAR(20);default:'pending';index"`
	ApprovalStatus  string `json:"approval_status" gorm:"type:VARCHAR(20);default:'pending';index"` // pending, approved, rejected
	ApprovedAt      string `json:"approved_at"`
	ApprovedBy      uint   `json:"approved_by"`
	RejectedBy      uint   `json:"rejected_by"`
	gorm.Model
}
