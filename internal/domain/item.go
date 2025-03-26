package domain

import (
	"gorm.io/gorm"
)

type Item struct {
	gorm.Model
	ID          int `gorm:"primaryKey"`
	Name        string
	Description string
	Quantity    int
}
