package db

import "github.com/google/uuid"

type ID string

func NewID() ID {
	return ID(NewIDString())
}

func NewIDString() string {
	return uuid.Must(uuid.NewV7()).String()
}
