package inbox

import (
	"context"

	"github.com/google/uuid"

	"payment-service/internal/domain/messages"
)

type ProcessorRepo interface {
	GetUnprocessed(ctx context.Context, limit int) ([]messages.InboxMessage, error)
	LockByID(ctx context.Context, id uuid.UUID) error
	MarkProcessed(ctx context.Context, id uuid.UUID) error
}
