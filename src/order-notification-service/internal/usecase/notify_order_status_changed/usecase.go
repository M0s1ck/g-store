package notify_order_status_changed

import (
	"context"
	"log"
)

type NotifyUsecase struct {
	notifier TransportNotifier
}

func NewUsecase(notifier TransportNotifier) *NotifyUsecase {
	return &NotifyUsecase{
		notifier: notifier,
	}
}

func (u *NotifyUsecase) Execute(ctx context.Context, evt Event) {
	log.Printf("Got event: %v %v %v", evt.Status, evt.OrderID, evt.UserID)

	// TODO: maybe add check if ok transition here (including registers of statuses)

	u.notifier.NotifyOrderStatusChanged(ctx, evt)
}
