package entity

import "errors"

var (
	ErrIDRequired   = errors.New("ID is required and must be valid")
	ErrNameRequired = errors.New("name is required")
	ErrNameTooLong  = errors.New("name cannot exceed 255 characters")
)

var (
	ErrInvalidPrice = errors.New("product price must be greater than zero")
)

var (
	ErrEmailRequired    = errors.New("email is required")
	ErrEmailTooLong     = errors.New("email cannot exceed 255 characters")
	ErrPasswordRequired = errors.New("password is required")
	ErrPasswordTooLong  = errors.New("password cannot exceed 255 characters")
)
