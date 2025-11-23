package entity

import (
	"github.com/jb-oliveira/fullcycle/tree/main/APIS/pkg/entity"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       entity.ID `json:"id" gorm:"type:uuid;primarykey"`
	Name     string    `json:"name" gorm:"column:usr_name;size:255"`
	Email    string    `json:"email" gorm:"column:usr_email;size:255;unique"`
	Password string    `json:"-" gorm:"column:usr_password;size:255"`
	entity.BaseModel
}

func (User) TableName() string {
	return "users"
}

func (u *User) Validate() error {
	if u.ID.String() == "" {
		return ErrIDRequired
	}
	if _, err := entity.ParseID(u.ID.String()); err != nil {
		return ErrIDRequired
	}
	if u.Name == "" {
		return ErrNameRequired
	}
	if len(u.Name) > 255 {
		return ErrNameTooLong
	}
	if u.Email == "" {
		return ErrEmailRequired
	}
	if len(u.Email) > 255 {
		return ErrEmailTooLong
	}
	if u.Password == "" {
		return ErrPasswordRequired
	}
	if len(u.Password) > 255 {
		return ErrPasswordTooLong
	}
	return nil
}

func (u *User) ValidatePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func NewUser(name, email, password string) (*User, error) {
	// Valida a password antes do hashing
	if password == "" {
		return nil, ErrPasswordRequired
	}
	if len(password) > 255 {
		return nil, ErrPasswordTooLong
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	user := &User{
		ID:       entity.NewID(),
		Name:     name,
		Email:    email,
		Password: string(hash),
	}
	if err := user.Validate(); err != nil {
		return nil, err
	}
	return user, nil
}
