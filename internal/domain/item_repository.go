package domain

type ItemRepository interface {
	FindAll() ([]Item, error)
	FindByID(id int) (Item, error)
	Create(item Item) (Item, error)
	Update(item Item) (Item, error)
	Delete(id int) error
}