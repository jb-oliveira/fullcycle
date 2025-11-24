package dto

// CreateProductInput
// usando o binding para quando for migrar pro gin
type CreateProductInput struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required,gt=0"`
}

// UpdateProductInput
type UpdateProductInput struct {
	Name  string  `json:"name" binding:"required"`
	Price float64 `json:"price" binding:"required,gt=0"`
}

// ProductOutput
type ProductOutput struct {
	ID    string  `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

// ProductListOutput
type ProductListOutput struct {
	Products []ProductOutput `json:"products"`
	Page     int             `json:"page"`
	Limit    int             `json:"limit"`
	Sort     string          `json:"sort"`
	Total    int             `json:"total"`
}
