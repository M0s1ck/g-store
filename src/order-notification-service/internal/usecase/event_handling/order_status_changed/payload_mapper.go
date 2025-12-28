package order_status_changed

type PayloadMapper interface {
	OrderStatusChangedEventFromPayload(payload []byte) (*Event, error)
}
