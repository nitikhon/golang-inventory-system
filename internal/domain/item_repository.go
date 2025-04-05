package domain

type ItemRepository interface {
	GetAllItems() ([]Item, error)
	GetItemByID(id int) (Item, error)
	Create(item Item) (Item, error)
	Update(item Item) (Item, error)
	Delete(id int) error
}