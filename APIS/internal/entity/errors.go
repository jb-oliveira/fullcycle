package entity

import "errors"

var (
	ErrIDRequired   = errors.New("ID é obrigatório e deve ser válido")
	ErrNameRequired = errors.New("nome é obrigatório")
	ErrNameTooLong  = errors.New("nome não pode exceder 255 caracteres")
)

var (
	ErrInvalidPrice = errors.New("preço do produto deve ser maior que zero")
)

var (
	ErrEmailRequired    = errors.New("email é obrigatório")
	ErrEmailTooLong     = errors.New("email não pode exceder 255 caracteres")
	ErrPasswordRequired = errors.New("senha é obrigatória")
	ErrPasswordTooLong  = errors.New("senha não pode exceder 255 caracteres")
)
