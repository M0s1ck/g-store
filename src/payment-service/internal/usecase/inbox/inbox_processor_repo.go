package inbox

import (
	"context"

	"github.com/google/uuid"

	"payment-service/internal/domain/messages"
)

type ProcessorRepo interface {
	MarkProcessed(ctx context.Context, id uuid.UUID) error
	GetUnprocessed(ctx context.Context, limit int) ([]messages.InboxMessage, error)
}
