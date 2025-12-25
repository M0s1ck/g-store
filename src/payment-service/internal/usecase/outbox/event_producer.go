package outbox

import (
	"context"
	"payment-service/internal/domain/messages"
)

type EventProducer interface {
	Publish(ctx context.Context, message *messages.OutboxMessage) error
}
