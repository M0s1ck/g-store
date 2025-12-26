package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"orders-service/internal/domain/entities"
	derrors "orders-service/internal/domain/errors"
	"orders-service/internal/infrastructure/db/postgres"
)

type OrderRepository struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) GetById(ctx context.Context, id uuid.UUID) (*entities.Order, error) {
	var order entities.Order
	err := r.db.GetContext(ctx, &order, "SELECT * FROM orders WHERE id = $1", id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: id=%s", derrors.ErrOrderNotFound, id)
	}

	if err != nil {
		return nil, err
	}

	return &order, err
}

func (r *OrderRepository) GetByUserId(ctx context.Context, userId uuid.UUID, page, limit int) ([]entities.Order, int, error) {
	var orders []entities.Order
	err := r.db.SelectContext(ctx, &orders,
		"SELECT * FROM orders WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3", userId, limit, (page-1)*limit)

	if err != nil {
		return nil, 0, err
	}

	var total int
	err = r.db.GetContext(ctx, &total,
		`SELECT COUNT(*) FROM orders WHERE user_id = $1`,
		userId,
	)
	if err != nil {
		return nil, 0, err
	}

	return orders, total, nil
}

func (r *OrderRepository) Create(ctx context.Context, order *entities.Order) error {
	exec := r.getExec(ctx)

	rows, err := exec.ExecContext(ctx,
		`INSERT INTO orders (id, user_id, amount, status)
				VALUES ($1, $2, $3, $4)`,
		order.Id, order.UserId, order.Amount, order.Status)

	log.Printf("%v rows effected", rows)

	return err
}

func (r *OrderRepository) UpdateStatus(ctx context.Context, order *entities.Order) error {
	exec := r.getExec(ctx)
	_, err := exec.ExecContext(ctx,
		"UPDATE orders SET status = $1, description = $2 WHERE id = $3",
		order.Status, order.Description, order.Id)
	return err
}

// returns sqlx.TX if we're in transaction or r.db if not
func (r *OrderRepository) getExec(ctx context.Context) sqlx.ExtContext {
	if tx, ok := ctx.Value(postgres.TxKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return r.db
}
