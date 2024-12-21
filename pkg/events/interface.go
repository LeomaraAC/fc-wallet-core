package events

import "time"

type EventInterface interface {
	GetName() string
	GetDateTime() time.Time
	GetPayload() interface{}
	SetPayload(payload interface{})
}

type EventHandlerInterface interface {
	Handle(event EventInterface)
}

type EventDispatcherInterface interface {
	Register(eventName string, Handler EventHandlerInterface) error
	Dispatch(event EventInterface) error
	Remove(eventName string, Handler EventHandlerInterface) error
	Has(eventName string, Handler EventHandlerInterface) bool
	Clear() error
}
