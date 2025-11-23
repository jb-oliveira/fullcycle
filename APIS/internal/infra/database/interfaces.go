package database

import "github.com/jb-oliveira/fullcycle/tree/main/APIS/internal/entity"

type UserInterface interface {
	Create(user *entity.User) error
	FindByEmail(email string) (*entity.User, error)
}
