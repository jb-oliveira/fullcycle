package database

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/jb-oliveira/fullcycle/APIS/internal/entity"
)

type User struct {
	db *gorm.DB
}

func NewUserDB(db *gorm.DB) *User {
	return &User{db: db}
}

func (u *User) Create(user *entity.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	return gorm.G[entity.User](u.db).Create(ctx, user)
}

func (u *User) FindByEmail(email string) (*entity.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*10)
	defer cancel()
	user, err := gorm.G[entity.User](u.db).Where("usr_email = ?", email).First(ctx)
	return &user, err
}
