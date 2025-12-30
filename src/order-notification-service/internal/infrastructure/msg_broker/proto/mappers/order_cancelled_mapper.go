package proto_mappers

import (
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"order-notification-service/internal/domain/events/consumed"
	"order-notification-service/internal/infrastructure/msg_broker/proto/gen"
)

type OrderCancelledPayloadMapper struct {
}

func NewOrderCancelledPayloadMapper() *OrderCancelledPayloadMapper {
	return &OrderCancelledPayloadMapper{}
}

func (m *OrderCancelledPayloadMapper) ToOrderCreatedEvent(payload []byte) (*consumed_events.OrderCancelledEvent, error) {
	var protoEvent gen.OrderCancelledEvent

	if err := proto.Unmarshal(payload, &protoEvent); err != nil {
		return nil, err
	}

	return protoToOrderCancelledEvent(&protoEvent)
}

func protoToOrderCancelledEvent(p *gen.OrderCancelledEvent) (*consumed_events.OrderCancelledEvent, error) {
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

	return &consumed_events.OrderCancelledEvent{
		OrderId:      oid,
		UserId:       uid,
		CancelReason: protoToCancellationReason(p.CancellationReason),
		CancelSource: protoToCancelSource(p.CancelSource),
		OccurredAt:   occurred,
	}, nil
}

func (m *OrderCancelledPayloadMapper) ToOrderCancelledEvent(payload []byte) (*consumed_events.OrderCancelledEvent, error) {
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

func protoToCancelSource(s gen.CancelSource) consumed_events.CancelSource {
	switch s {
	case gen.CancelSource_CANCEL_SOURCE_STORE:
		return consumed_events.CancelSourceStore
	case gen.CancelSource_CANCEL_SOURCE_CUSTOMER:
		return consumed_events.CancelSourceCustomer
	default:
		return ""
	}
}

func protoToCancellationReason(r gen.CancellationReason) consumed_events.OrderCancellationReason {
	switch r {
	case gen.CancellationReason_NO_PAYMENT_ACCOUNT:
		return consumed_events.CancellationNoPaymentAccount
	case gen.CancellationReason_INSUFFICIENT_FUNDS:
		return consumed_events.CancellationInsufficientFunds
	case gen.CancellationReason_PAYMENT_INTERNAL_ERROR:
		return consumed_events.CancellationPaymentInternalError
	case gen.CancellationReason_OUT_OF_STOCK:
		return consumed_events.CancellationOutOfStock
	case gen.CancellationReason_DELIVERY_UNAVAILABLE:
		return consumed_events.CancellationDeliveryUnavailable
	case gen.CancellationReason_CHANGED_MIND:
		return consumed_events.CancellationChangedMind
	default:
		return ""
	}
}
