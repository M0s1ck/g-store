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

func (u *NotifyUsecase) Execute(ctx context.Context, req Request) {
	log.Printf("Got request: %v %v %v", req.Status, req.OrderID, req.UserID)
}
