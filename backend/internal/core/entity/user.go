package entity

type User struct {
	GormModel
	Username     string `json:"username" gorm:"index:idx_username"`
	Email        string `json:"email" gorm:"index:idx_email"`
	Password     string `json:"password"`
	Phone        string `json:"phone" gorm:"index:idx_phone"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	IsAdmin      bool   `json:"is_admin" gorm:"default:false"`
	RefreshToken string `json:"refresh_token"`
}

type Token struct {
    AccessToken  string
    RefreshToken string
}

type PublicUser struct {
    ID        uint   `json:"id"`
    Username  string `json:"username"`
    FirstName string `json:"first_name"`
    LastName  string `json:"last_name"`
}

// for GORM table recognition
func (PublicUser) TableName() string {
    return "users"
}