package entity

import (
	"errors"

	"github.com/jb-oliveira/fullcycle/tree/main/APIS/pkg/entity"
)

var (
	// ErrIDRequired is returned when product ID is empty or invalid.
	ErrIDRequired = errors.New("product ID is required and must be valid")
	// ErrNameRequired is returned when product name is empty.
	ErrNameRequired = errors.New("product name is required")
	// ErrInvalidPrice is returned when product price is less than or equal to zero.
	ErrInvalidPrice = errors.New("product price must be greater than zero")
)

type Product struct {
	entity.IDModel
	Name  string  `json:"name" gorm:"column:prd_name"`
	Price float64 `json:"price" gorm:"column:prd_price"`
}

// Validate checks if the Product fields are valid.
func (p *Product) Validate() error {
	if p.ID.String() == "" {
		return ErrIDRequired
	}
	if _, err := entity.ParseID(p.ID.String()); err != nil {
		return ErrIDRequired
	}
	if p.Name == "" {
		return ErrNameRequired
	}
	if p.Price <= 0 {
		return ErrInvalidPrice
	}
	return nil
}

// NewProduct creates a new Product with the given name and price.
// It generates a new ID and validates the product before returning.
func NewProduct(name string, price float64) (*Product, error) {
	product := &Product{
		IDModel: entity.IDModel{
			ID: entity.NewID(),
		},
		Name:  name,
		Price: price,
	}
	if err := product.Validate(); err != nil {
		return nil, err
	}
	return product, nil
}
