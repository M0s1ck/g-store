package repository

import (
	"context"
	"log"

	"github.com/jmoiron/sqlx"

	"orders-service/internal/infrastructure/db/postgres"
	"orders-service/internal/usecase/common/outbox"
)

type OutboxRepository struct {
	db *sqlx.DB
}

func NewOutboxRepository(db *sqlx.DB) *OutboxRepository {
	return &OutboxRepository{db: db}
}

func (r *OutboxRepository) Create(ctx context.Context, model *outbox.Model) error {
	exec := r.getExec(ctx)

	result, err := exec.ExecContext(ctx, `
        INSERT INTO outbox (
            id,
            aggregate,
            aggregate_id,
            event_type,
            payload
        ) VALUES ($1, $2, $3, $4, $5)
    `,
		model.Id,
		model.Aggregate,
		model.AggregateID,
		model.EventType,
		model.Payload,
	)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	log.Printf("%d rows affected", rows)

	return nil
}

// returns sqlx.TX if we're in transaction or r.db if not
func (r *OutboxRepository) getExec(ctx context.Context) sqlx.ExtContext {
	if tx, ok := ctx.Value(postgres.TxKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return r.db
}
