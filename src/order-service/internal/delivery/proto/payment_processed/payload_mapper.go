package payment_processed

import (
	"fmt"
	"orders-service/internal/usecase/event_handlers/payment_processed"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type PayloadMapper struct {
}

func NewPayloadMapper() *PayloadMapper {
	return &PayloadMapper{}
}

func (p *PayloadMapper) ToPaymentProcessedEvent(payload []byte) (*payment_processed.PaymentProcessedEvent, error) {
	var protoEvt PaymentProcessedEvent

	if err := proto.Unmarshal(payload, &protoEvt); err != nil {
		return nil, err
	}

	event, err := paymentProcessedEventFromProto(&protoEvt)
	if err != nil {
		return nil, err
	}

	return event, nil
}

func paymentProcessedEventFromProto(
	prEvt *PaymentProcessedEvent,
) (*payment_processed.PaymentProcessedEvent, error) {

	// --- UUIDs ---
	messageId, err := uuid.Parse(prEvt.MessageId)
	if err != nil {
		return nil, fmt.Errorf("invalid message_id: %w", err)
	}

	orderId, err := uuid.Parse(prEvt.OrderId)
	if err != nil {
		return nil, fmt.Errorf("invalid order_id: %w", err)
	}

	userId, err := uuid.Parse(prEvt.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	// --- status ---
	var status payment_processed.PaymentStatus
	switch prEvt.Status {
	case PaymentStatus_PAYMENT_STATUS_SUCCESS:
		status = payment_processed.PaymentSuccess

	case PaymentStatus_PAYMENT_STATUS_FAILED:
		status = payment_processed.PaymentFailed

	default:
		return nil, fmt.Errorf("unknown payment status: %v", prEvt.Status)
	}

	// --- nullable failure reason ---
	var failureReason *payment_processed.PaymentFailureReason
	if prEvt.PaymentFailureReason != nil {
		var fr payment_processed.PaymentFailureReason

		switch *prEvt.PaymentFailureReason {
		case PaymentFailureReason_PAYMENT_FAILURE_REASON_NO_ACCOUNT:
			fr = payment_processed.FailureNoAccount

		case PaymentFailureReason_PAYMENT_FAILURE_REASON_INSUFFICIENT_FUNDS:
			fr = payment_processed.FailureInsufficientFunds

		case PaymentFailureReason_PAYMENT_FAILURE_REASON_INTERNAL_ERROR:
			fr = payment_processed.FailureInternal

		case PaymentFailureReason_PAYMENT_FAILURE_REASON_UNSPECIFIED:
			fr = ""

		default:
			return nil, fmt.Errorf(
				"unknown payment failure reason: %v",
				*prEvt.PaymentFailureReason,
			)
		}

		if fr != "" {
			failureReason = &fr
		}
	}

	// --- occurred at ---
	var occurredAt time.Time
	if prEvt.OccurredAt != nil {
		occurredAt = prEvt.OccurredAt.AsTime()
	}

	return &payment_processed.PaymentProcessedEvent{
		MessageId:            messageId,
		OrderId:              orderId,
		UserId:               userId,
		Amount:               prEvt.Amount,
		Status:               status,
		PaymentFailureReason: failureReason,
		OccurredAt:           occurredAt,
	}, nil
}
