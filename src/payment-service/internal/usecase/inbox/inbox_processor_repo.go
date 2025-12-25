package inbox

import (
	"context"
	"payment-service/internal/domain/messages"

	"github.com/google/uuid"
)

type ProcessorRepo interface {
	GetUnprocessed(ctx context.Context, limit int) ([]messages.InboxMessage, error)
	LockByID(ctx context.Context, id uuid.UUID) error
	MarkProcessed(ctx context.Context, id uuid.UUID) error
}
