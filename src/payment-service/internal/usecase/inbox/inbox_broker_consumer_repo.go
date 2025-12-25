package inbox

import (
	"context"

	"payment-service/internal/domain/messages"
)

// BrokerConsumerRepo used by broker consumer
type BrokerConsumerRepo interface {
	SaveIdempotent(ctx context.Context, msg messages.InboxMessage) error
}
