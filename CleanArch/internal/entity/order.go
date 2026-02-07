package entity

import "errors"

type Order struct {
	ID    string
	Price float64
	Tax   float64
}

func NewOrder(id string, price float64, tax float64) (*Order, error) {
	order := &Order{
		ID:    id,
		Price: price,
		Tax:   tax,
	}
	err := order.Validate()
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (o *Order) Validate() error {
	if o.ID == "" {
		return errors.New("ID is required")
	}
	if o.Price <= 0 {
		return errors.New("Price must be greater than 0")
	}
	if o.Tax <= 0 {
		return errors.New("Tax must be greater than 0")
	}
	return nil
}

func (o *Order) CalculateFinalPrice() float64 {
	return o.Price + o.Tax
}
