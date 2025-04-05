package entity

import (
	"gorm.io/gorm"
)

type Borrowing struct {
	UserID           uint `gorm:"primaryKey"`
	ItemID           uint `gorm:"primaryKey"`
	BorrowedAt       string
	ReturnedAt       string
	DueDate          string
	borrowedQuantity int
	borrowingStatus  string `gorm:"type:VARCHAR(20);default:'active'"`  // active, returned, overdue
	approvalStatus   string `gorm:"type:VARCHAR(20);default:'pending'"` // pending, approved, rejected
	approvedAt       string
	approvedBy       uint
	gorm.Model
}
