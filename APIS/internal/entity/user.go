package entity

import (
	"github.com/jb-oliveira/fullcycle/tree/main/APIS/pkg/entity"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	entity.IDModel
	Name     string `json:"name" gorm:"column:usr_name"`
	Email    string `json:"email" gorm:"column:usr_email"`
	Password string `json:"-" gorm:"column:usr_password"`
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func NewUser(name, email, password string) (*User, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	return &User{
		IDModel: entity.IDModel{
			ID: entity.NewID(),
		},
		Name:     name,
		Email:    email,
		Password: string(hash),
	}, nil
}
