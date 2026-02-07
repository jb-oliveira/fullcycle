package event

import (
	"sync"
	"time"
)

type Event interface {
	GetID() string
	GetName() string
	GetPayload() any
	GetCreatedAt() time.Time
	SetPayload(payload any) error
}

type EventHandler interface {
	Handle(event Event, wg *sync.WaitGroup) error
}

type EventDispatcher interface {
	Register(eventName string, handler EventHandler) error
	Dispatch(event Event) error
	Remove(eventName string, handler EventHandler) error
}
