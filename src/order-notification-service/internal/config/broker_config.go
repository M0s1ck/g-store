package config

type BrokerConfig struct {
	Brokers                     []string
	ConsumerGroup               string
	OrderEventTopic             string
	OrderStatusChangedEventType string
	OrderCancelledEventType     string
	AllowedEventTypes           map[string]struct{}
}
