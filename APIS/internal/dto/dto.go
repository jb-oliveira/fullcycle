package dto

type ErrorResponse struct {
	Messages []string `json:"messages"`
	Code     int      `json:"code"`
}

type AuthResponse struct {
	Token string `json:"token"`
}

// CreateProductInput
// using binding for when migrating to gin
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

type CreateUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserOutput struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}
