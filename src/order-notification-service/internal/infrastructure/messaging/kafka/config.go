package kafka

import "order-notification-service/internal/config"

type Config struct {
	Brokers                     []string
	ConsumerGroup               string
	OrderNotificationEventTopic string
	OrderStatusChangedEventType string
}

func NewKafkaConfig(conf *config.BrokerConfig) *Config {
	return &Config{
		Brokers:                     conf.Brokers,
		ConsumerGroup:               conf.ConsumerGroup,
		OrderNotificationEventTopic: conf.OrderNotificationEventTopic,
		OrderStatusChangedEventType: conf.OrderStatusChangedEventType,
	}
}
