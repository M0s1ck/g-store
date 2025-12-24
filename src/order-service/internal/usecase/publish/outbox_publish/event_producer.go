package outbox_publish

import (
	"context"

	"orders-service/internal/domain/messages"
)

type EventProducer interface {
	Publish(ctx context.Context, message *messages.OutboxMessage) error
}
