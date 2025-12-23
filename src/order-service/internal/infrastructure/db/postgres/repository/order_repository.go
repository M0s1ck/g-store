package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"

	"orders-service/internal/domain/entities"
	derrors "orders-service/internal/domain/errors"
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
