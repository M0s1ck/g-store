package order_created

import (
	"context"
	"fmt"
	"log"
	"payment-service/internal/domain/entities"
	"payment-service/internal/domain/errors"

	"payment-service/internal/domain/messages"

	"github.com/google/uuid"
)

type OrderCreatedEventHandler struct {
	orderCreatedEventType string
	mapper                OrderCreatedEventMapper
	bTransactionRepo      BalanceTransactionRepoCreator
	accountRepo           AccountRepoBalanceUpdater
}

func NewOrderCreatedEventHandler(
	accountRepo AccountRepoBalanceUpdater,
	bTransactionRepo BalanceTransactionRepoCreator,
	mapper OrderCreatedEventMapper,
	orderCreatedEventType string,
) *OrderCreatedEventHandler {

	return &OrderCreatedEventHandler{
		accountRepo:           accountRepo,
		bTransactionRepo:      bTransactionRepo,
		orderCreatedEventType: orderCreatedEventType,
		mapper:                mapper,
	}
}

func (o OrderCreatedEventHandler) EventType() string {
	return o.orderCreatedEventType
}

// Handle updates account balances, adds balance transaction, saves ... event message to outbox.
// Should be called in outer transaction with inbox-set-processed
func (o OrderCreatedEventHandler) Handle(ctx context.Context, msg messages.InboxMessage) error {
	log.Printf("Got msg!!! : %v %v", msg.Topic, msg.Payload)

	event, err := o.mapper.ToOrderCreatedEvent(msg.Payload)
	if err != nil {
		return err
	}

	// TODO: add check id=f err is not found -> means we processed? (to outbox)
	account, err := o.accountRepo.GetByUserId(ctx, event.UserId)
	if err != nil {
		return err
	}

	// TODO: at this point we save to outbox
	if account.Balance < event.Amount {
		return fmt.Errorf("account %v: %w", account.ID, errors.ErrInsufficientFunds)
	}

	account.Balance -= event.Amount

	err = o.accountRepo.UpdateBalance(ctx, account)
	if err != nil {
		return err
	}

	balanceTransaction := entities.BalanceTransaction{
		ID:        uuid.New(),
		AccountID: account.ID,
		OrderID:   &event.OrderId,
		Amount:    event.Amount,
		Type:      entities.TransactionPayment,
	}

	// TODO: maybe at this point we save to outbox
	err = o.bTransactionRepo.Create(ctx, &balanceTransaction)
	if err != nil {
		return err
	}

	//TODO: potential issue with fk

	// TODO: maybe at this point we save to outbox

	return nil
}
