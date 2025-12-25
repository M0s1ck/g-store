package order_created

import (
	"context"
	"errors"
	"log"
	"payment-service/internal/domain/entities"
	myerrors "payment-service/internal/domain/errors"
	"payment-service/internal/domain/events"
	"payment-service/internal/usecase/common/outbox"
	"payment-service/internal/usecase/order_created/repo_contracts"
	"time"

	"payment-service/internal/domain/messages"

	"github.com/google/uuid"
)

type OrderCreatedEventHandler struct {
	accountRepo           repo_contracts.AccountRepoBalanceUpdater
	bTransactionRepo      repo_contracts.BalanceTransactionRepoCreator
	outboxRepo            outbox.RepositoryCreator
	outboxMsgFactory      *outbox.MessageFactory
	orderMapper           OrderCreatedEventMapper
	orderCreatedEventType string
}

func NewOrderCreatedEventHandler(
	accountRepo repo_contracts.AccountRepoBalanceUpdater,
	bTransactionRepo repo_contracts.BalanceTransactionRepoCreator,
	outboxRepo outbox.RepositoryCreator,
	outboxMsgFactory *outbox.MessageFactory,
	mapper OrderCreatedEventMapper,
	orderCreatedEventType string,
) *OrderCreatedEventHandler {

	return &OrderCreatedEventHandler{
		accountRepo:           accountRepo,
		bTransactionRepo:      bTransactionRepo,
		outboxRepo:            outboxRepo,
		outboxMsgFactory:      outboxMsgFactory,
		orderCreatedEventType: orderCreatedEventType,
		orderMapper:           mapper,
	}
}

func (o *OrderCreatedEventHandler) EventType() string {
	return o.orderCreatedEventType
}

// Handle updates account balances, adds balance transaction, saves payment-processed event message to outbox.
// Should be called in outer transaction with inbox-set-processed
func (o *OrderCreatedEventHandler) Handle(ctx context.Context, msg messages.InboxMessage) error {
	log.Printf("Got msg!!! : %v %v", msg.Topic, msg.Payload)

	orderEvent, err := o.orderMapper.ToOrderCreatedEvent(msg.Payload)
	if err != nil {
		return err
	}

	account, err := o.accountRepo.GetByUserId(ctx, orderEvent.UserId)

	if errors.Is(err, myerrors.ErrAccountNotFound) {
		return o.handlePaymentFailure(ctx, orderEvent, events.FailureNoAccount)
	}

	if err != nil {
		log.Printf("error getting account by usId %v: %s", orderEvent.UserId, err)
		return err
	}

	if account.Balance < orderEvent.Amount {
		return o.handlePaymentFailure(ctx, orderEvent, events.FailureInsufficientFunds)
	}

	account.Balance -= orderEvent.Amount
	err = o.accountRepo.UpdateBalance(ctx, account)
	if err != nil {
		log.Printf("error updating balance %v: %s", account.ID, err)
		return err
	}

	balanceTransaction := createBalanceTransaction(orderEvent, account)

	err = o.bTransactionRepo.Create(ctx, balanceTransaction)
	if err != nil {
		log.Printf("error creating balance-transaction amount=%v: %s", balanceTransaction.Amount, err)
		return err
	}

	paymentEvent := createPaymentProcessedEvent(orderEvent, events.PaymentSuccess, nil)
	outboxMsg := o.outboxMsgFactory.PaymentProcessedEventToOutboxMessage(paymentEvent)
	return o.outboxRepo.Create(ctx, outboxMsg)
}

func (o *OrderCreatedEventHandler) handlePaymentFailure(
	ctx context.Context,
	orderEvent *OrderCreatedEvent,
	reason events.PaymentFailureReason) error {

	paymentEvent := createPaymentProcessedEvent(orderEvent, events.PaymentFailed, &reason)
	outboxMsg := o.outboxMsgFactory.PaymentProcessedEventToOutboxMessage(paymentEvent)
	return o.outboxRepo.Create(ctx, outboxMsg)
}

func createPaymentProcessedEvent(
	orderCreated *OrderCreatedEvent,
	status events.PaymentStatus,
	fail *events.PaymentFailureReason,
) *events.PaymentProcessedEvent {

	return &events.PaymentProcessedEvent{
		MessageId:            uuid.New(),
		OrderId:              orderCreated.OrderId,
		UserId:               orderCreated.UserId,
		Amount:               orderCreated.Amount,
		Status:               status,
		PaymentFailureReason: fail,
		OccurredAt:           time.Now(),
	}
}

func createBalanceTransaction(orderEvent *OrderCreatedEvent, account *entities.Account) *entities.BalanceTransaction {
	return &entities.BalanceTransaction{
		ID:        uuid.New(),
		AccountID: account.ID,
		OrderID:   &orderEvent.OrderId,
		Amount:    orderEvent.Amount,
		Type:      entities.TransactionPayment,
	}
}
