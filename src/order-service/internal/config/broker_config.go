package config

type BrokerConfig struct {
	Brokers                   []string
	ConsumerGroup             string
	OrderEventsTopic          string
	OrderCreatedEventType     string
	PaymentEventsTopic        string
	PaymentProcessedEventType string
}
