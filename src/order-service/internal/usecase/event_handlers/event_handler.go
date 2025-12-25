package event_handlers

import "context"

type EventHandler interface {
	EventType() string
	Handle(ctx context.Context, payload []byte) error
}
