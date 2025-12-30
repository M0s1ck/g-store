package pay_for_order

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"

	"payment-service/internal/domain/entities"
	myerrors "payment-service/internal/domain/errors"
	"payment-service/internal/domain/events/produced"
	"payment-service/internal/usecase/common"
	"payment-service/internal/usecase/common/outbox"
)

type PayUsecase struct {
	accRepo          AccountRepoBalanceUpdater
	txRepo           BalanceTransactionRepoCreator
	txManager        common.TxManager
	outboxRepo       outbox.RepositoryCreator
	outboxMsgFactory *outbox.MessageFactory
}

func NewPayUsecase(
	accRepo AccountRepoBalanceUpdater,
	txRepo BalanceTransactionRepoCreator,
	txManager common.TxManager,
	outboxRepo outbox.RepositoryCreator,
	outboxMsgFactory *outbox.MessageFactory,
) *PayUsecase {

	return &PayUsecase{
		accRepo:          accRepo,
		txRepo:           txRepo,
		txManager:        txManager,
		outboxRepo:       outboxRepo,
		outboxMsgFactory: outboxMsgFactory,
	}
}

func (uc *PayUsecase) Execute(ctx context.Context, cmd *Command) error {
	return uc.txManager.WithinTx(ctx, func(ctx context.Context) error {

		account, err := uc.accRepo.GetByUserId(ctx, cmd.UserId)

		if errors.Is(err, myerrors.ErrAccountNotFound) {
			return uc.handlePaymentFailure(ctx, cmd, produced_events.FailureNoAccount)
		}

		if err != nil {
			log.Printf("error getting account by usId %v: %s", cmd.UserId, err)
			return err
		}

		if account.Balance < cmd.Amount {
			return uc.handlePaymentFailure(ctx, cmd, produced_events.FailureInsufficientFunds)
		}

		account.Balance -= cmd.Amount
		err = uc.accRepo.UpdateBalance(ctx, account)
		if err != nil {
			log.Printf("error updating balance %v: %s", account.ID, err)
			return err
		}

		balanceTransaction := createBalanceTransaction(cmd, account)

		err = uc.txRepo.Create(ctx, balanceTransaction)
		if err != nil {
			log.Printf("error creating balance-transaction amount=%v: %s", balanceTransaction.Amount, err)
			return err
		}

		paymentEvent := createPaymentProcessedEvent(cmd, produced_events.PaymentSuccess, nil)
		outboxMsg := uc.outboxMsgFactory.PaymentProcessedEventToOutboxMessage(paymentEvent)
		return uc.outboxRepo.Create(ctx, outboxMsg)
	})
}

func (uc *PayUsecase) handlePaymentFailure(
	ctx context.Context,
	cmd *Command,
	reason produced_events.PaymentFailureReason) error {

	paymentEvent := createPaymentProcessedEvent(cmd, produced_events.PaymentFailed, &reason)
	outboxMsg := uc.outboxMsgFactory.PaymentProcessedEventToOutboxMessage(paymentEvent)
	return uc.outboxRepo.Create(ctx, outboxMsg)
}

func createPaymentProcessedEvent(
	cmd *Command,
	status produced_events.PaymentStatus,
	fail *produced_events.PaymentFailureReason,
) *produced_events.PaymentProcessedEvent {

	return &produced_events.PaymentProcessedEvent{
		OrderId:              cmd.OrderId,
		UserId:               cmd.UserId,
		Amount:               cmd.Amount,
		Status:               status,
		PaymentFailureReason: fail,
		OccurredAt:           time.Now(),
	}
}

func createBalanceTransaction(cmd *Command, account *entities.Account) *entities.BalanceTransaction {
	return &entities.BalanceTransaction{
		ID:        uuid.New(),
		AccountID: account.ID,
		OrderID:   &cmd.OrderId,
		Amount:    cmd.Amount,
		Type:      entities.TransactionPayment,
		Direction: entities.DirectionOut,
	}
}
