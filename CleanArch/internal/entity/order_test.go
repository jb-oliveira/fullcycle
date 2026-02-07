package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewOrder_Success(t *testing.T) {
	order, err := NewOrder("123", 10.0, 2.0)
	assert.Nil(t, err)
	assert.Equal(t, "123", order.ID)
	assert.Equal(t, 10.0, order.Price)
	assert.Equal(t, 2.0, order.Tax)
}

func TestNewOrder_EmptyID(t *testing.T) {
	order, err := NewOrder("", 10.0, 2.0)
	assert.NotNil(t, err)
	assert.Equal(t, "ID is required", err.Error())
	assert.Nil(t, order)
}

func TestNewOrder_InvalidPrice(t *testing.T) {
	order, err := NewOrder("123", 0, 2.0)
	assert.NotNil(t, err)
	assert.Equal(t, "Price must be greater than 0", err.Error())
	assert.Nil(t, order)
}

func TestNewOrder_InvalidTax(t *testing.T) {
	order, err := NewOrder("123", 10.0, 0)
	assert.NotNil(t, err)
	assert.Equal(t, "Tax must be greater than 0", err.Error())
	assert.Nil(t, order)
}

func TestOrder_CalculateFinalPrice(t *testing.T) {
	order, err := NewOrder("123", 10.0, 2.0)
	assert.Nil(t, err)
	finalPrice, err := order.CalculateFinalPrice()
	assert.Nil(t, err)
	assert.Equal(t, 12.0, finalPrice)
	assert.Equal(t, 12.0, order.FinalPrice)
}
