package entity

import (
	"testing"

	"github.com/jb-oliveira/fullcycle/tree/main/APIS/pkg/entity"
	"github.com/stretchr/testify/assert"
)

// TestNewProduct tests the creation of a new product with valid data.
// It verifies that a product is created successfully with all fields populated correctly.
func TestNewProduct(t *testing.T) {
	product, err := NewProduct("Laptop", 999.99)

	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.NotEmpty(t, product.ID)
	assert.Equal(t, "Laptop", product.Name)
	assert.Equal(t, 999.99, product.Price)
}

// TestNewProduct_ValidatesFields tests the validation logic in NewProduct.
// It uses table-driven tests to verify that validation errors are returned
// for invalid inputs (empty name, zero price, negative price).
func TestNewProduct_ValidatesFields(t *testing.T) {
	tests := []struct {
		name        string
		productName string
		price       float64
		expectError error
	}{
		{
			name:        "valid product",
			productName: "Mouse",
			price:       29.99,
			expectError: nil,
		},
		{
			name:        "empty name",
			productName: "",
			price:       50.0,
			expectError: ErrNameRequired,
		},
		{
			name:        "zero price",
			productName: "Keyboard",
			price:       0,
			expectError: ErrInvalidPrice,
		},
		{
			name:        "negative price",
			productName: "Monitor",
			price:       -100.0,
			expectError: ErrInvalidPrice,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			product, err := NewProduct(tt.productName, tt.price)

			if tt.expectError != nil {
				assert.ErrorIs(t, err, tt.expectError)
				assert.Nil(t, product)
			} else {
				assert.Nil(t, err)
				assert.NotNil(t, product)
				assert.Equal(t, tt.productName, product.Name)
				assert.Equal(t, tt.price, product.Price)
			}
		})
	}
}

// TestNewProduct_GeneratesUniqueIDs verifies that each product created
// receives a unique ID, ensuring no ID collisions occur.
func TestNewProduct_GeneratesUniqueIDs(t *testing.T) {
	product1, err1 := NewProduct("Product One", 10.0)
	product2, err2 := NewProduct("Product Two", 20.0)

	assert.Nil(t, err1)
	assert.Nil(t, err2)
	assert.NotEqual(t, product1.ID, product2.ID, "Each product should have a unique ID")
}

// TestProduct_Validate tests the Validate method directly on Product instances.
// It verifies that validation rules are correctly enforced for various
// product configurations.
func TestProduct_Validate(t *testing.T) {
	tests := []struct {
		name        string
		product     *Product
		expectError error
	}{
		{
			name: "valid product",
			product: &Product{
				IDModel: entity.IDModel{
					ID: entity.NewID(),
				},
				Name:  "Valid Product",
				Price: 100.0,
			},
			expectError: nil,
		},
		{
			name: "empty name",
			product: &Product{
				IDModel: entity.IDModel{
					ID: entity.NewID(),
				},
				Name:  "",
				Price: 50.0,
			},
			expectError: ErrNameRequired,
		},
		{
			name: "zero price",
			product: &Product{
				IDModel: entity.IDModel{
					ID: entity.NewID(),
				},
				Name:  "Product",
				Price: 0,
			},
			expectError: ErrInvalidPrice,
		},
		{
			name: "negative price",
			product: &Product{
				IDModel: entity.IDModel{
					ID: entity.NewID(),
				},
				Name:  "Product",
				Price: -10.0,
			},
			expectError: ErrInvalidPrice,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.product.Validate()

			if tt.expectError != nil {
				assert.ErrorIs(t, err, tt.expectError)
			} else {
				assert.Nil(t, err)
			}
		})
	}
}

// TestNewProduct_WithSpecialCharactersInName verifies that product names
// can contain special characters without causing validation errors.
func TestNewProduct_WithSpecialCharactersInName(t *testing.T) {
	product, err := NewProduct("Product @#$% Special!", 99.99)

	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, "Product @#$% Special!", product.Name)
}

// TestNewProduct_WithVerySmallPrice tests that very small prices (like 0.01)
// are accepted as valid, ensuring the validation allows minimum price values.
func TestNewProduct_WithVerySmallPrice(t *testing.T) {
	product, err := NewProduct("Cheap Item", 0.01)

	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, 0.01, product.Price)
}

// TestNewProduct_WithVeryLargePrice tests that very large prices are accepted,
// ensuring there's no upper limit validation that would reject valid high prices.
func TestNewProduct_WithVeryLargePrice(t *testing.T) {
	product, err := NewProduct("Expensive Item", 999999.99)

	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, 999999.99, product.Price)
}

// TestNewProduct_WithLongName verifies that product names can be very long
// without causing validation errors or truncation.
func TestNewProduct_WithLongName(t *testing.T) {
	longName := "This is a very long product name that contains many characters and should still be valid"
	product, err := NewProduct(longName, 50.0)

	assert.Nil(t, err)
	assert.NotNil(t, product)
	assert.Equal(t, longName, product.Name)
}

// TestProduct_Validate_WithZeroUUID tests the edge case where a product
// has a zero UUID. Currently, the implementation accepts this as valid.
func TestProduct_Validate_WithZeroUUID(t *testing.T) {
	product := &Product{
		IDModel: entity.IDModel{
			ID: entity.ID{},
		},
		Name:  "Product",
		Price: 100.0,
	}

	err := product.Validate()
	assert.Nil(t, err)
}
