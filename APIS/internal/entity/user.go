package entity

import (
	"errors"

	"github.com/jb-oliveira/fullcycle/tree/main/APIS/pkg/entity"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserIDRequired    = errors.New("ID do usuário é obrigatório e deve ser válido")
	ErrUserNameRequired  = errors.New("nome do usuário é obrigatório")
	ErrUserNameTooLong   = errors.New("nome do usuário não pode exceder 255 caracteres")
	ErrUserEmailRequired = errors.New("email do usuário é obrigatório")
	ErrUserEmailTooLong  = errors.New("email do usuário não pode exceder 255 caracteres")
	ErrPasswordRequired  = errors.New("senha do usuário é obrigatória")
	ErrPasswordTooLong   = errors.New("senha do usuário não pode exceder 255 caracteres")
)

type User struct {
	entity.IDModel
	Name     string `json:"name" gorm:"column:usr_name;size:255"`
	Email    string `json:"email" gorm:"column:usr_email;size:255"`
	Password string `json:"-" gorm:"column:usr_password;size:255"`
}

func (User) TableName() string {
	return "users"
}

func (u *User) Validate() error {
	if u.ID.String() == "" {
		return ErrUserIDRequired
	}
	if _, err := entity.ParseID(u.ID.String()); err != nil {
		return ErrUserIDRequired
	}
	if u.Name == "" {
		return ErrUserNameRequired
	}
	if len(u.Name) > 255 {
		return ErrUserNameTooLong
	}
	if u.Email == "" {
		return ErrUserEmailRequired
	}
	if len(u.Email) > 255 {
		return ErrUserEmailTooLong
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
	// Validate password before hashing
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
		IDModel: entity.IDModel{
			ID: entity.NewID(),
		},
		Name:     name,
		Email:    email,
		Password: string(hash),
	}
	if err := user.Validate(); err != nil {
		return nil, err
	}
	return user, nil
}
