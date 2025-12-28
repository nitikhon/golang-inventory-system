package entity


type Item struct {
	GormModel
	Name            string `json:"name" gorm:"uniqueIndex:idx_name"`
	Description     string `json:"description"`
	AvailableAmount int    `json:"available_amount"`
	TotalAmount     int    `json:"total_amount" gorm:"not null;default:0"`
	Status          string `json:"status" gorm:"type:VARCHAR(20);default:'available'"` // available, borrowed, maintenance, lost
}
