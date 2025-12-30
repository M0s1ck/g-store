package refund_order

import (
	"context"

	"github.com/google/uuid"

	"payment-service/internal/domain/entities"
)

type BalanceTransactionRepoCreator interface {
	Create(ctx context.Context, transaction *entities.BalanceTransaction) error
	ExistsRefundByOrderId(ctx context.Context, orderId uuid.UUID) (bool, error)
	GetPaymentByOrderId(ctx context.Context, orderId uuid.UUID) (*entities.BalanceTransaction, error)
}
