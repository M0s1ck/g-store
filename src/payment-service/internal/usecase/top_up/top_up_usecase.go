package top_up

import (
	"context"

	"github.com/google/uuid"

	"payment-service/internal/domain/entities"
	myerrors "payment-service/internal/domain/errors"
	"payment-service/internal/usecase/common"
)

type TopUpUsecase struct {
	txManger common.TxManager
	accRepo  AccountRepoBalanceUpdater
	bTxRepo  BalanceTransactionRepoCreator
}

func NewTopUpUsecase(
	accRepo AccountRepoBalanceUpdater,
	bTxRepo BalanceTransactionRepoCreator,
	txManger common.TxManager) *TopUpUsecase {

	return &TopUpUsecase{
		accRepo:  accRepo,
		bTxRepo:  bTxRepo,
		txManger: txManger,
	}
}

func (uc *TopUpUsecase) Execute(ctx context.Context,
	accId uuid.UUID,
	userId uuid.UUID,
	amount int64) error {

	if amount <= 0 {
		return myerrors.ErrAmountNotPositive
	}

	return uc.txManger.WithinTx(ctx, func(ctx context.Context) error {

		account, err := uc.accRepo.GetById(ctx, accId)
		if err != nil {
			return err
		}

		if account.UserID != userId {
			return myerrors.ErrForbidden
		}

		account.Balance += amount

		err = uc.accRepo.UpdateBalance(ctx, account)
		if err != nil {
			return err
		}

		balanceTx := entities.BalanceTransaction{
			ID:        uuid.New(),
			AccountID: accId,
			Amount:    amount,
			Type:      entities.TransactionTopUp,
			Direction: entities.DirectionIn,
		}

		err = uc.bTxRepo.Create(ctx, &balanceTx)
		return err
	})
}
