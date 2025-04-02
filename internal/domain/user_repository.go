package domain

type UserRepository interface {
	CreateUser(user *User) (*User, error)
	UpdateUser(user *User) (*User, error)
	DeleteUser(id int) error
	GetAllUsers() ([]*User, error)
	GetUserByID(id int) (*User, error)
	GetUserByUsername(username string) (*User, error)
	GetUserByEmail(email string) (*User, error)
	GetUserByPhone(phone string) (*User, error)
}
