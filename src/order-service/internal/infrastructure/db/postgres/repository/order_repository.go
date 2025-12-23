package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"orders-service/internal/domain/entities"
	derrors "orders-service/internal/domain/errors"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
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
