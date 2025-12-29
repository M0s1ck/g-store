package proto_mappers

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
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

// TODO: remove from here

// PayloadToOrderCancelledEvent deserializes payload bytes into domain event.
func (p OrderCancelledPayloadMapper) PayloadToOrderCancelledEvent(payload []byte) (*published_events.OrderCancelledEvent, error) {
	var protoEvt gen.OrderCancelledEvent
	if err := proto.Unmarshal(payload, &protoEvt); err != nil {
		log.Println("Error while unmarshalling payload to proto OrderCancelledEvent:", err)
		return nil, err
	}

	domainEvt, err := protoToOrderCancelledEvent(&protoEvt)
	if err != nil {
		log.Println("Error while converting proto OrderCancelledEvent to domain:", err)
		return nil, err
	}

	return domainEvt, nil
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

func protoToOrderCancelledEvent(p *gen.OrderCancelledEvent) (*published_events.OrderCancelledEvent, error) {
	if p == nil {
		return nil, errors.New("nil proto OrderCancelledEvent")
	}

	oid, err := uuid.Parse(p.OrderId)
	if err != nil {
		return nil, err
	}

	uid, err := uuid.Parse(p.UserId)
	if err != nil {
		return nil, err
	}

	var occurred time.Time
	if p.OccurredAt != nil {
		occurred = p.OccurredAt.AsTime()
	} else {
		occurred = time.Time{}
	}

	return &published_events.OrderCancelledEvent{
		OrderId:      oid,
		UserId:       uid,
		CancelReason: protoToCancellationReason(p.CancellationReason),
		CancelSource: protoToCancelSource(p.CancelSource),
		OccurredAt:   occurred,
	}, nil
}

func protoToCancelSource(s gen.CancelSource) published_events.CancelSource {
	switch s {
	case gen.CancelSource_CANCEL_SOURCE_STORE:
		return published_events.CancelSourceStore
	case gen.CancelSource_CANCEL_SOURCE_CUSTOMER:
		return published_events.CancelSourceCustomer
	default:
		return ""
	}
}

func protoToCancellationReason(r gen.CancellationReason) value_objects.CancellationReason {
	switch r {
	case gen.CancellationReason_NO_PAYMENT_ACCOUNT:
		return value_objects.CancellationNoPaymentAccount
	case gen.CancellationReason_INSUFFICIENT_FUNDS:
		return value_objects.CancellationInsufficientFunds
	case gen.CancellationReason_PAYMENT_INTERNAL_ERROR:
		return value_objects.CancellationPaymentInternalError
	case gen.CancellationReason_OUT_OF_STOCK:
		return value_objects.CancellationOutOfStock
	case gen.CancellationReason_DELIVERY_UNAVAILABLE:
		return value_objects.CancellationDeliveryUnavailable
	case gen.CancellationReason_CHANGED_MIND:
		return value_objects.CancellationChangedMind
	default:
		return ""
	}
}
