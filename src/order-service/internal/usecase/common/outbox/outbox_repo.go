package outbox

import "context"

type Repository interface {
	Create(ctx context.Context, model *Model) error
}
