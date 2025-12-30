package proto_mappers

import (
	"log"

	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	"orders-service/internal/domain/events/produced"
	"orders-service/internal/domain/value_objects"
	"orders-service/internal/infrastructure/msg_broker/proto/gen"
)

type OrderCancelledPayloadMapper struct {
}

func NewOrderCancelledPayloadMapper() *OrderCancelledPayloadMapper {
	return &OrderCancelledPayloadMapper{}
}

// OrderCancelledEventToPayload serializes domain event to protobuf bytes.
func (p OrderCancelledPayloadMapper) OrderCancelledEventToPayload(event *published_events.OrderCancelledEvent) ([]byte, error) {
	protoEvt := orderCancelledEventToProto(event)

	payload, err := proto.Marshal(protoEvt)
	if err != nil {
		log.Println("Error while marshalling proto order-cancelled-event to payload:", err)
		return nil, err
	}

	return payload, nil
}

func orderCancelledEventToProto(event *published_events.OrderCancelledEvent) *gen.OrderCancelledEvent {
	if event == nil {
		return &gen.OrderCancelledEvent{}
	}

	return &gen.OrderCancelledEvent{
		OrderId:            event.OrderId.String(),
		UserId:             event.UserId.String(),
		CancellationReason: cancellationReasonToProto(event.CancelReason),
		CancelSource:       cancelSourceToProto(event.CancelSource),
		OccurredAt:         timestamppb.New(event.OccurredAt),
	}
}

func cancelSourceToProto(s published_events.CancelSource) gen.CancelSource {
	switch s {
	case published_events.CancelSourceStore:
		return gen.CancelSource_CANCEL_SOURCE_STORE
	case published_events.CancelSourceCustomer:
		return gen.CancelSource_CANCEL_SOURCE_CUSTOMER
	default:
		return gen.CancelSource_CANCEL_SOURCE_UNSPECIFIED
	}
}

func cancellationReasonToProto(r value_objects.CancellationReason) gen.CancellationReason {
	switch r {
	case value_objects.CancellationNoPaymentAccount:
		return gen.CancellationReason_NO_PAYMENT_ACCOUNT
	case value_objects.CancellationInsufficientFunds:
		return gen.CancellationReason_INSUFFICIENT_FUNDS
	case value_objects.CancellationPaymentInternalError:
		return gen.CancellationReason_PAYMENT_INTERNAL_ERROR
	case value_objects.CancellationOutOfStock:
		return gen.CancellationReason_OUT_OF_STOCK
	case value_objects.CancellationDeliveryUnavailable:
		return gen.CancellationReason_DELIVERY_UNAVAILABLE
	case value_objects.CancellationChangedMind:
		return gen.CancellationReason_CHANGED_MIND
	default:
		return gen.CancellationReason_CANCELLATION_REASON_UNSPECIFIED
	}
}
