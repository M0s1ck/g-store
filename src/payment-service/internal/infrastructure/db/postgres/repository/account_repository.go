package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"payment-service/internal/domain/entities"
	myerrors "payment-service/internal/domain/errors"
	"payment-service/internal/infrastructure/db/postgres"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type AccountRepository struct {
	db *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) GetById(ctx context.Context, id uuid.UUID) (*entities.Account, error) {
	exec := r.getExec(ctx)
	var account entities.Account
	err := sqlx.GetContext(ctx, exec, &account, "SELECT * FROM accounts WHERE id = $1", id)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: id=%s", myerrors.ErrAccountNotFound, id)
	}

	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) GetByUserId(ctx context.Context, userId uuid.UUID) (*entities.Account, error) {
	exec := r.getExec(ctx)
	var account entities.Account
	err := sqlx.GetContext(ctx, exec, &account, "SELECT * FROM accounts WHERE user_id = $1", userId)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("%w: for user_id=%s", myerrors.ErrAccountNotFound, userId)
	}

	if err != nil {
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepository) UpdateBalance(ctx context.Context, account *entities.Account) error {
	exec := r.getExec(ctx)
	res, err := exec.ExecContext(ctx,
		"UPDATE accounts SET balance = $1 WHERE id = $2 ",
		account.Balance, account.ID)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return fmt.Errorf("account %s: %w", account.ID, myerrors.ErrAccountNotFound)
	}

	return nil
}

// returns sqlx.TX if we're in transaction or r.db if not
func (r *AccountRepository) getExec(ctx context.Context) sqlx.ExtContext {
	if tx, ok := ctx.Value(postgres.TxKey{}).(*sqlx.Tx); ok {
		return tx
	}
	return r.db
}
