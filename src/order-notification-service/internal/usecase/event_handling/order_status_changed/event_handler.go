package order_status_changed

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
	ordStChangedEvtType string,
) *EventHandler {

	return &EventHandler{
		notifyUC:            notifyUC,
		ordStChangedEvtType: ordStChangedEvtType,
		payloadMapper:       mapper,
	}
}

func (h *EventHandler) EventType() string {
	return h.ordStChangedEvtType
}

func (h *EventHandler) Handle(ctx context.Context, payload []byte) error {
	event, err := h.payloadMapper.OrderStatusChangedEventFromPayload(payload)
	if err != nil {
		return err
	}

	ucReq := notify_order_status_changed.Request{
		OrderID:            event.OrderId,
		UserID:             event.UserId,
		Status:             event.Status,
		CancellationReason: event.CancellationReason,
	}

	h.notifyUC.Execute(ctx, ucReq)

	return nil
}
