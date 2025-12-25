package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"payment-service/internal/domain/messages"
	"payment-service/internal/infrastructure/db/postgres"
)

type InboxRepository struct {
	db *sqlx.DB
}

func NewInboxRepository(db *sqlx.DB) *InboxRepository {
	return &InboxRepository{db: db}
}

func (r *InboxRepository) SaveIdempotent(ctx context.Context, msg messages.InboxMessage) error {
	exec := r.getExec(ctx)

	_, err := exec.ExecContext(ctx, `
		INSERT INTO inbox (id, topic, event_type, key, payload)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id) DO NOTHING
	`, msg.Id, msg.Topic, msg.EventType, msg.Key, msg.Payload)

	return err
}

func (r *InboxRepository) GetUnprocessed(ctx context.Context, limit int) ([]messages.InboxMessage, error) {
	var msgs []messages.InboxMessage

	exec := r.getExec(ctx)

	err := sqlx.SelectContext(ctx, exec, &msgs, `
		SELECT * FROM inbox
		WHERE processed_at IS NULL
		ORDER BY created_at
		LIMIT $1
	`, limit)

	return msgs, err
}

func (r *InboxRepository) MarkProcessed(ctx context.Context, id uuid.UUID) error {
	exec := r.getExec(ctx)
	_, err := exec.ExecContext(ctx,
		`UPDATE inbox SET processed_at = now() WHERE id = $1`, id)
	return err
}

func (r *InboxRepository) LockByID(ctx context.Context, id uuid.UUID) error {
	exec := r.getExec(ctx)
	_, err := exec.ExecContext(ctx, "SELECT * FROM inbox WHERE id = $1 FOR UPDATE SKIP LOCKED", id)
	return err
}

// returns sqlx.TX if we're in transaction or r.db if not
func (r *InboxRepository) getExec(ctx context.Context) sqlx.ExtContext {
	if tx, ok := ctx.Value(postgres.TxKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return r.db
}
