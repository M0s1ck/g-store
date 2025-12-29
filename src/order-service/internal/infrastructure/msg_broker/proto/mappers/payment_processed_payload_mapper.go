package proto_mappers

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"orders-service/internal/domain/events/consumed"
	"orders-service/internal/infrastructure/msg_broker/proto/gen"
)

type PaymentProcessedPayloadMapper struct {
}

func NewPaymentProcessedPayloadMapper() *PaymentProcessedPayloadMapper {
	return &PaymentProcessedPayloadMapper{}
}

func (p *PaymentProcessedPayloadMapper) ToPaymentProcessedEvent(payload []byte) (*consumed_events.PaymentProcessedEvent, error) {
	var protoEvt gen.PaymentProcessedEvent

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
	prEvt *gen.PaymentProcessedEvent,
) (*consumed_events.PaymentProcessedEvent, error) {

	orderId, err := uuid.Parse(prEvt.OrderId)
	if err != nil {
		return nil, fmt.Errorf("invalid order_id: %w", err)
	}

	userId, err := uuid.Parse(prEvt.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	// --- status ---
	var status consumed_events.PaymentStatus
	switch prEvt.Status {
	case gen.PaymentStatus_PAYMENT_STATUS_SUCCESS:
		status = consumed_events.PaymentSuccess

	case gen.PaymentStatus_PAYMENT_STATUS_FAILED:
		status = consumed_events.PaymentFailed

	default:
		return nil, fmt.Errorf("unknown payment status: %v", prEvt.Status)
	}

	// --- nullable failure reason ---
	var failureReason *consumed_events.PaymentFailureReason
	if prEvt.PaymentFailureReason != nil {
		var fr consumed_events.PaymentFailureReason

		switch *prEvt.PaymentFailureReason {
		case gen.PaymentFailureReason_PAYMENT_FAILURE_REASON_NO_ACCOUNT:
			fr = consumed_events.FailureNoAccount

		case gen.PaymentFailureReason_PAYMENT_FAILURE_REASON_INSUFFICIENT_FUNDS:
			fr = consumed_events.FailureInsufficientFunds

		case gen.PaymentFailureReason_PAYMENT_FAILURE_REASON_INTERNAL_ERROR:
			fr = consumed_events.FailureInternal

		case gen.PaymentFailureReason_PAYMENT_FAILURE_REASON_UNSPECIFIED:
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

	return &consumed_events.PaymentProcessedEvent{
		OrderId:              orderId,
		UserId:               userId,
		Amount:               prEvt.Amount,
		Status:               status,
		PaymentFailureReason: failureReason,
		OccurredAt:           occurredAt,
	}, nil
}
