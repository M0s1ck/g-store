package proto_mappers

import (
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"payment-service/internal/domain/events/produced"
	"payment-service/internal/infrastructure/msg_broker/proto/gen"
)

type PaymentProcessedPayloadMapper struct {
}

func NewPaymentProcessedPayloadMapper() *PaymentProcessedPayloadMapper {
	return &PaymentProcessedPayloadMapper{}
}

func (mapper *PaymentProcessedPayloadMapper) PaymentProcessedEventToPayload(event *produced_events.PaymentProcessedEvent) ([]byte, error) {
	eventProto := paymentProcessedEventToProto(event)

	payload, err := proto.Marshal(eventProto)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func paymentProcessedEventToProto(
	event *produced_events.PaymentProcessedEvent,
) *gen.PaymentProcessedEvent {

	// --- status ---
	var status gen.PaymentStatus
	switch event.Status {
	case produced_events.PaymentSuccess:
		status = gen.PaymentStatus_PAYMENT_STATUS_SUCCESS
	case produced_events.PaymentFailed:
		status = gen.PaymentStatus_PAYMENT_STATUS_FAILED
	default:
		status = gen.PaymentStatus_PAYMENT_STATUS_UNSPECIFIED
	}

	// --- nullable failure reason ---
	var failureReason *gen.PaymentFailureReason
	if event.PaymentFailureReason != nil {
		var fr gen.PaymentFailureReason

		switch *event.PaymentFailureReason {
		case produced_events.FailureNoAccount:
			fr = gen.PaymentFailureReason_PAYMENT_FAILURE_REASON_NO_ACCOUNT
		case produced_events.FailureInsufficientFunds:
			fr = gen.PaymentFailureReason_PAYMENT_FAILURE_REASON_INSUFFICIENT_FUNDS
		case produced_events.FailureInternal:
			fr = gen.PaymentFailureReason_PAYMENT_FAILURE_REASON_INTERNAL_ERROR
		default:
			fr = gen.PaymentFailureReason_PAYMENT_FAILURE_REASON_UNSPECIFIED
		}

		failureReason = &fr
	}

	// --- occurred at ---
	var occurredAt *timestamppb.Timestamp
	if !event.OccurredAt.IsZero() {
		occurredAt = timestamppb.New(event.OccurredAt)
	}

	return &gen.PaymentProcessedEvent{
		OrderId:              event.OrderId.String(),
		UserId:               event.UserId.String(),
		Amount:               event.Amount,
		Status:               status,
		PaymentFailureReason: failureReason,
		OccurredAt:           occurredAt,
	}
}
