package entity

import (
	"gorm.io/gorm"
)

type User struct {
	Username  string `json:"username" gorm:"index:idx_username"`
	Email     string `json:"email" gorm:"index:idx_email"`
	Password  string `json:"password"`
	Phone     string `json:"phone" gorm:"index:idx_phone"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	IsAdmin   bool   `json:"is_admin" gorm:"default:false"`
	gorm.Model
}
