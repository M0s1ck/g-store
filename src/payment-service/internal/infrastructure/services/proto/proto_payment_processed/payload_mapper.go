package proto_payment_processed

import (
	"payment-service/internal/domain/events"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PayloadMapper struct {
}

func NewPayloadMapper() *PayloadMapper {
	return &PayloadMapper{}
}

func (mapper *PayloadMapper) PaymentProcessedEventToPayload(event *events.PaymentProcessedEvent) ([]byte, error) {
	eventProto := paymentProcessedEventToProto(event)

	payload, err := proto.Marshal(eventProto)
	if err != nil {
		return nil, err
	}

	return payload, nil
}

func paymentProcessedEventToProto(
	event *events.PaymentProcessedEvent,
) *PaymentProcessedEvent {

	// --- status ---
	var status PaymentStatus
	switch event.Status {
	case events.PaymentSuccess:
		status = PaymentStatus_PAYMENT_STATUS_SUCCESS
	case events.PaymentFailed:
		status = PaymentStatus_PAYMENT_STATUS_FAILED
	default:
		status = PaymentStatus_PAYMENT_STATUS_UNSPECIFIED
	}

	// --- nullable failure reason ---
	var failureReason *PaymentFailureReason
	if event.PaymentFailureReason != nil {
		var fr PaymentFailureReason

		switch *event.PaymentFailureReason {
		case events.FailureNoAccount:
			fr = PaymentFailureReason_PAYMENT_FAILURE_REASON_NO_ACCOUNT
		case events.FailureInsufficientFunds:
			fr = PaymentFailureReason_PAYMENT_FAILURE_REASON_INSUFFICIENT_FUNDS
		case events.FailureInternal:
			fr = PaymentFailureReason_PAYMENT_FAILURE_REASON_INTERNAL_ERROR
		default:
			fr = PaymentFailureReason_PAYMENT_FAILURE_REASON_UNSPECIFIED
		}

		failureReason = &fr
	}

	// --- occurred at ---
	var occurredAt *timestamppb.Timestamp
	if !event.OccurredAt.IsZero() {
		occurredAt = timestamppb.New(event.OccurredAt)
	}

	return &PaymentProcessedEvent{
		MessageId:            event.MessageId.String(),
		OrderId:              event.OrderId.String(),
		UserId:               event.UserId.String(),
		Amount:               event.Amount,
		Status:               status,
		PaymentFailureReason: failureReason,
		OccurredAt:           occurredAt,
	}
}
