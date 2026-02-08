package entity

type OrderRepository interface {
	Save(order *Order) error
	FindAll(page, limit int, sort, sortDir string) ([]Order, error)
}
