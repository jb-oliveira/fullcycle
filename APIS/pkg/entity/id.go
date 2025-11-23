package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ID = uuid.UUID

type IDModel struct {
	// ID        string `gorm:"primarykey;default:uuidv7()::varchar"`
	ID        ID `gorm:"type:uuid;primarykey;default:uuidv7()" json:"id"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func NewID() ID {
	return ID(uuid.Must(uuid.NewV7()))
}

func ParseID(s string) (ID, error) {
	uuid, err := uuid.Parse(s)
	return ID(uuid), err
}
