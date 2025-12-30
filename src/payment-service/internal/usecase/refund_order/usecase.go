package refund_order

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"

	"payment-service/internal/domain/entities"
	"payment-service/internal/usecase/common"
)

type RefundUsecase struct {
	accountRepo AccountRepoBalanceUpdater
	txRepo      BalanceTransactionRepoCreator
	txManager   common.TxManager
}

func NewRefundUsecase(
	accRepo AccountRepoBalanceUpdater,
	txRepo BalanceTransactionRepoCreator,
	txManager common.TxManager,
) *RefundUsecase {

	return &RefundUsecase{
		accountRepo: accRepo,
		txRepo:      txRepo,
		txManager:   txManager,
	}
}

// Execute should be called within tx
func (uc *RefundUsecase) Execute(ctx context.Context, cmd *Command) error {
	return uc.txManager.WithinTx(ctx, func(ctx context.Context) error {

		exists, err := uc.txRepo.ExistsRefundByOrderId(ctx, cmd.OrderId)
		if err != nil {
			return err
		}

		if exists {
			return nil
		}

		payment, err := uc.txRepo.GetPaymentByOrderId(ctx, cmd.OrderId)
		if err != nil {
			log.Printf("err getting payment by order id=%s: %s", cmd.OrderId, err)
			return err
		}

		account, err := uc.accountRepo.GetById(ctx, payment.AccountID)
		if err != nil {
			return err
		}

		refund := refundFromPayment(payment)
		account.Balance += refund.Amount

		err = uc.txRepo.Create(ctx, refund)
		if err != nil {
			return err
		}

		return uc.accountRepo.UpdateBalance(ctx, account)
	})

}

func refundFromPayment(pay *entities.BalanceTransaction) *entities.BalanceTransaction {
	return &entities.BalanceTransaction{
		ID:        uuid.New(),
		AccountID: pay.AccountID,
		OrderID:   pay.OrderID,
		Amount:    pay.Amount,
		Direction: entities.DirectionIn,
		Type:      entities.TransactionRefund,
		CreatedAt: time.Now(),
	}
}
