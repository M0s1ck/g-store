package inbox

import (
	"context"

	"payment-service/internal/domain/messages"
)

// MessageHandler interface to different event handlers
type MessageHandler interface {
	EventType() string
	Handle(ctx context.Context, msg messages.InboxMessage) error
}
