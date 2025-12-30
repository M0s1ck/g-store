package config

type BrokerConfig struct {
	Brokers                   []string
	ConsumerGroup             string
	OrderEventsTopic          string
	OrderCreatedEventType     string
	OrderCancelledEventType   string
	PaymentEventsTopic        string
	PaymentProcessedEventType string
}
