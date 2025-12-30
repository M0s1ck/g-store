package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"payment-service/internal/domain/entities"
	"payment-service/internal/infrastructure/db/postgres"
)

type BalanceTransactionRepository struct {
	db *sqlx.DB
}

func NewBalanceTransactionRepository(db *sqlx.DB) *BalanceTransactionRepository {
	return &BalanceTransactionRepository{db: db}
}

func (r *BalanceTransactionRepository) Create(ctx context.Context, transaction *entities.BalanceTransaction) error {
	exec := r.getExec(ctx)
	_, err := exec.ExecContext(ctx,
		"INSERT INTO balance_transactions(id, account_id, amount, order_id, type, direction) VALUES ($1, $2, $3, $4, $5, $6)",
		transaction.ID, transaction.AccountID, transaction.Amount, transaction.OrderID, transaction.Type, transaction.Direction,
	)

	return err
}

func (r *BalanceTransactionRepository) ExistsRefundByOrderId(ctx context.Context, orderId uuid.UUID) (bool, error) {
	const q = `SELECT EXISTS (SELECT 1 FROM balance_transactions WHERE order_id = $1 AND type = 'REFUND')`

	var exists bool
	err := r.db.QueryRowContext(ctx, q, orderId).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *BalanceTransactionRepository) GetPaymentByOrderId(ctx context.Context, orderId uuid.UUID) (*entities.BalanceTransaction, error) {
	const q = `SELECT * FROM balance_transactions WHERE order_id = $1 AND type = 'PAYMENT'`

	var tx entities.BalanceTransaction
	err := r.db.GetContext(ctx, &tx, q, orderId)
	if err != nil {
		return nil, err
	}

	return &tx, nil
}

// returns sqlx.TX if we're in transaction or r.db if not
func (r *BalanceTransactionRepository) getExec(ctx context.Context) sqlx.ExtContext {
	if tx, ok := ctx.Value(postgres.TxKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return r.db
}
