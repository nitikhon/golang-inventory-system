package entity

import (
	"gorm.io/gorm"
)

type Item struct {
	Name        string `gorm:"uniqueIndex:idx_name"`
	Description string
	AvailableAmount   int
	TotalAmount int    `gorm:"not null;default:0"`
	Status      string `gorm:"type:VARCHAR(20);default:'available'"` // available, borrowed, maintenance, lost
	gorm.Model
}
