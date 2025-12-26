package repository

import (
	"context"

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
		"INSERT INTO balance_transactions(id, account_id, amount, order_id, type) VALUES ($1, $2, $3, $4, $5)",
		transaction.ID, transaction.AccountID, transaction.Amount, transaction.OrderID, transaction.Type,
	)

	return err
}

// returns sqlx.TX if we're in transaction or r.db if not
func (r *BalanceTransactionRepository) getExec(ctx context.Context) sqlx.ExtContext {
	if tx, ok := ctx.Value(postgres.TxKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return r.db
}
