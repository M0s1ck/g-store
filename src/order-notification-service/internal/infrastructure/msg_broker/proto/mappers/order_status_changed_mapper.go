package proto_mappers

import (
	"fmt"
	"order-notification-service/internal/domain/events/consumed"
	"order-notification-service/internal/infrastructure/msg_broker/proto/gen"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type OrderStatusChangedPayloadMapper struct {
}

func NewOrderStatusChangedPayloadMapper() *OrderStatusChangedPayloadMapper {
	return &OrderStatusChangedPayloadMapper{}
}

func (p *OrderStatusChangedPayloadMapper) OrderStatusChangedEventFromPayload(payload []byte) (*consumed_events.OrderStatusChangedEvent, error) {
	var evtProto gen.OrderStatusChangedEvent

	err := proto.Unmarshal(payload, &evtProto)
	if err != nil {
		return nil, err
	}

	return orderStatusChangedEventFromProto(&evtProto)
}

func orderStatusChangedEventFromProto(
	prEvt *gen.OrderStatusChangedEvent,
) (*consumed_events.OrderStatusChangedEvent, error) {

	orderId, err := uuid.Parse(prEvt.OrderId)
	if err != nil {
		return nil, fmt.Errorf("invalid order_id: %w", err)
	}

	userId, err := uuid.Parse(prEvt.UserId)
	if err != nil {
		return nil, fmt.Errorf("invalid user_id: %w", err)
	}

	// --- status ---
	var status consumed_events.OrderStatus
	switch prEvt.Status {
	case gen.OrderStatus_ORDER_STATUS_PENDING:
		status = consumed_events.OrderPending
	case gen.OrderStatus_ORDER_STATUS_PAID:
		status = consumed_events.OrderPaid
	case gen.OrderStatus_ORDER_STATUS_CANCELLED:
		status = consumed_events.OrderCanceled
	case gen.OrderStatus_ORDER_STATUS_ASSEMBLING:
		status = consumed_events.OrderAssembling
	case gen.OrderStatus_ORDER_STATUS_ASSEMBLED:
		status = consumed_events.OrderAssembled
	case gen.OrderStatus_ORDER_STATUS_DELIVERING:
		status = consumed_events.OrderDelivering
	case gen.OrderStatus_ORDER_STATUS_ISSUED:
		status = consumed_events.OrderIssued
	default:
		return nil, fmt.Errorf("unknown order status: %v", prEvt.Status)
	}

	// --- created at ---
	var createdAt time.Time
	if prEvt.CreatedAt != nil {
		createdAt = prEvt.CreatedAt.AsTime()
	}

	return &consumed_events.OrderStatusChangedEvent{
		OrderId:    orderId,
		UserId:     userId,
		Status:     status,
		OccurredAt: createdAt,
	}, nil
}
