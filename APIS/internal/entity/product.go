package entity

import (
	"errors"

	"github.com/jb-oliveira/fullcycle/tree/main/APIS/pkg/entity"
)

var (
	ErrIDRequired   = errors.New("ID do produto é obrigatório e deve ser válido")
	ErrNameRequired = errors.New("nome do produto é obrigatório")
	ErrInvalidPrice = errors.New("preço do produto deve ser maior que zero")
)

type Product struct {
	entity.IDModel
	Name  string  `json:"name" gorm:"column:prd_name;size:255"`
	Price float64 `json:"price" gorm:"column:prd_price;type:decimal(10,2)"`
}

func (Product) TableName() string {
	return "products"
}

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
