package proto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"

	"order-notification-service/internal/usecase/domain/value_objects"
	"order-notification-service/internal/usecase/event_handling/order_status_changed"
)

type PayloadMapper struct {
}

func NewPayloadMapper() *PayloadMapper {
	return &PayloadMapper{}
}

func (p *PayloadMapper) OrderStatusChangedEventFromPayload(payload []byte) (*order_status_changed.Event, error) {
	var evtProto OrderStatusChangedEvent

	err := proto.Unmarshal(payload, &evtProto)
	if err != nil {
		return nil, err
	}

	return orderStatusChangedEventFromProto(&evtProto)
}

func orderStatusChangedEventFromProto(
	prEvt *OrderStatusChangedEvent,
) (*order_status_changed.Event, error) {

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
	var status value_objects.OrderStatus
	switch prEvt.Status {
	case OrderStatus_ORDER_STATUS_PENDING:
		status = value_objects.OrderPending
	case OrderStatus_ORDER_STATUS_PAID:
		status = value_objects.OrderPaid
	case OrderStatus_ORDER_STATUS_CANCELED:
		status = value_objects.OrderCanceled
	case OrderStatus_ORDER_STATUS_ASSEMBLING:
		status = value_objects.OrderAssembling
	case OrderStatus_ORDER_STATUS_ASSEMBLED:
		status = value_objects.OrderAssembled
	case OrderStatus_ORDER_STATUS_DELIVERING:
		status = value_objects.OrderDelivering
	case OrderStatus_ORDER_STATUS_ISSUED:
		status = value_objects.OrderIssued
	default:
		return nil, fmt.Errorf("unknown order status: %v", prEvt.Status)
	}

	// --- cancellation reason ---
	var cancellationReason *string
	if prEvt.CancellationReason != nil && *prEvt.CancellationReason != "" {
		cancellationReason = prEvt.CancellationReason
	}

	// --- created at ---
	var createdAt time.Time
	if prEvt.CreatedAt != nil {
		createdAt = prEvt.CreatedAt.AsTime()
	}

	return &order_status_changed.Event{
		MessageId:          messageId,
		OrderId:            orderId,
		UserId:             userId,
		Status:             status,
		CancellationReason: cancellationReason,
		CreatedAt:          createdAt,
	}, nil
}
