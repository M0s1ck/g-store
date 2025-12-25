package order_created

type OrderCreatedEventMapper interface {
	ToOrderCreatedEvent(payload []byte) (*OrderCreatedEvent, error)
}
