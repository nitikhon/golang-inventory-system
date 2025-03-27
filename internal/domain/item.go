package domain

import (
	"gorm.io/gorm"
)

type Item struct {
	ID          int `gorm:"primaryKey"`
	Name        string
	Description string
	Quantity    int
	gorm.Model
}
