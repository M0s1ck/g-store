package config

type BrokerConfig struct {
	Brokers                     []string
	ConsumerGroup               string
	OrderNotificationEventTopic string
	OrderStatusChangedEventType string
}
