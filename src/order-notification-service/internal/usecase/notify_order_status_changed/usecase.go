package notify_order_status_changed

import (
	"context"
	"log"
	consumed_events "order-notification-service/internal/domain/events/consumed"
)

type NotifyUsecase struct {
	notifier TransportNotifier
}

func NewUsecase(notifier TransportNotifier) *NotifyUsecase {
	return &NotifyUsecase{
		notifier: notifier,
	}
}

func (u *NotifyUsecase) Execute(ctx context.Context, evt consumed_events.OrderStatusChangedEvent) {
	log.Printf("Got event: %v %v %v", evt.Status, evt.OrderId, evt.UserId)

	// TODO: maybe add check if ok transition here (including registers of statuses)

	u.notifier.NotifyOrderStatusChanged(ctx, evt)
}
