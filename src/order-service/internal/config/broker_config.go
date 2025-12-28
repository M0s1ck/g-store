package config

type BrokerConfig struct {
	Brokers                     []string
	ConsumerGroup               string
	OrderCommandEventsTopic     string
	OrderCreatedEventType       string
	PaymentEventsTopic          string
	PaymentProcessedEventType   string
	OrderNotificationEventTopic string
	OrderStatusChangedEventType string
}
