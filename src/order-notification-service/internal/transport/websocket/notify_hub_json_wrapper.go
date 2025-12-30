package websocket

import (
	"context"
	"encoding/json"
	"log"
	"order-notification-service/internal/domain/events/consumed"
)

type NotifyHubJSONWrapper struct {
	hub *Hub
}

func NewNotifyHubJSONWrapper(hub *Hub) *NotifyHubJSONWrapper {
	return &NotifyHubJSONWrapper{
		hub: hub,
	}
}

func (n *NotifyHubJSONWrapper) NotifyOrderStatusChanged(
	_ context.Context,
	event *consumed_events.OrderStatusChangedEvent,
) {
	dto := StatusChanged{
		OrderID: event.OrderId,
		Status:  string(event.Status),
	}

	data, err := json.Marshal(dto)

	if err != nil {
		log.Printf("marshal error: %v", err)
		return
	}

	envelope := Envelope{
		Type: "order.status.changed",
		Data: data,
	}

	payload, err := json.Marshal(envelope)
	if err != nil {
		log.Printf("marshal envelope error: %v", err)
		return
	}

	n.hub.NotifyOrder(event.OrderId, payload)
}

func (n *NotifyHubJSONWrapper) NotifyOrderCancelled(
	_ context.Context,
	event *consumed_events.OrderCancelledEvent,
) {
	dto := OrderCancelled{
		OrderID: event.OrderId,
		Reason:  string(event.CancelReason),
	}

	data, err := json.Marshal(dto)

	if err != nil {
		log.Printf("marshal error: %v", err)
		return
	}

	envelope := Envelope{
		Type: "order.cancelled",
		Data: data,
	}

	payload, err := json.Marshal(envelope)
	if err != nil {
		log.Printf("marshal envelope error: %v", err)
		return
	}

	n.hub.NotifyOrder(event.OrderId, payload)
}
