package order_cancelled

import (
	"context"
	"order-notification-service/internal/usecase/notify_order_status_changed"
)

type EventHandler struct {
	notifyUC            *notify_order_status_changed.NotifyUsecase
	payloadMapper       PayloadMapper
	ordStChangedEvtType string
}

func NewEventHandler(
	notifyUC *notify_order_status_changed.NotifyUsecase,
	mapper PayloadMapper,
	ordCancelEvtType string,
) *EventHandler {

	return &EventHandler{
		notifyUC:            notifyUC,
		ordStChangedEvtType: ordCancelEvtType,
		payloadMapper:       mapper,
	}
}

func (h *EventHandler) EventType() string {
	return h.ordStChangedEvtType
}

func (h *EventHandler) Handle(ctx context.Context, payload []byte) error {
	event, err := h.payloadMapper.ToOrderCancelledEvent(payload)
	if err != nil {
		return err
	}

	h.notifyUC.NotifyCancelled(ctx, event)
	return nil
}
