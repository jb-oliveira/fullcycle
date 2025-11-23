package dto

// CreateProductInput represents the input data for creating a new product
type CreateProductInput struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required,gt=0"`
}

// UpdateProductInput represents the input data for updating an existing product
type UpdateProductInput struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required,gt=0"`
}

// ProductOutput represents the output data for a product
type ProductOutput struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// ProductListOutput represents a paginated list of products
type ProductListOutput struct {
	Products []ProductOutput `json:"products"`
	Page     int             `json:"page"`
	Limit    int             `json:"limit"`
	Total    int             `json:"total"`
}
