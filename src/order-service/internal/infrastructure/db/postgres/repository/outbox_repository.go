package repository

import (
	"context"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"orders-service/internal/domain/messages"
	"orders-service/internal/infrastructure/db/postgres"
)

type OutboxRepository struct {
	db *sqlx.DB
}

func NewOutboxRepository(db *sqlx.DB) *OutboxRepository {
	return &OutboxRepository{db: db}
}

func (r *OutboxRepository) Create(ctx context.Context, model *messages.OutboxMessage) error {
	exec := r.getExec(ctx)

	result, err := exec.ExecContext(ctx, `
        INSERT INTO outbox (
            id,
            aggregate,
            aggregate_id,
            event_type,
            key,
            payload
        ) VALUES ($1, $2, $3, $4, $5, $6)
    `,
		model.Id,
		model.Aggregate,
		model.AggregateID,
		model.EventType,
		model.Key,
		model.Payload,
	)
	if err != nil {
		return err
	}

	rows, _ := result.RowsAffected()
	log.Printf("%d rows affected", rows)

	return nil
}

func (r *OutboxRepository) GetUnsent(ctx context.Context, limit int) ([]messages.OutboxMessage, error) {
	var msgs []messages.OutboxMessage
	err := r.db.SelectContext(ctx, &msgs,
		"SELECT * FROM outbox WHERE sent_at IS NULL ORDER BY created_at LIMIT $1", limit)
	return msgs, err
}

func (r *OutboxRepository) MarkAsSent(ctx context.Context, id uuid.UUID) error {
	exec := r.getExec(ctx)
	_, err := exec.ExecContext(ctx,
		`UPDATE outbox SET sent_at = now() WHERE id = $1`, id)
	return err
}

func (r *OutboxRepository) IncrementRetry(ctx context.Context, id uuid.UUID) error {
	exec := r.getExec(ctx)
	_, err := exec.ExecContext(ctx,
		`UPDATE outbox SET retry_count = retry_count + 1 WHERE id = $1`, id)
	return err
}

// returns sqlx.TX if we're in transaction or r.db if not
func (r *OutboxRepository) getExec(ctx context.Context) sqlx.ExtContext {
	if tx, ok := ctx.Value(postgres.TxKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return r.db
}
