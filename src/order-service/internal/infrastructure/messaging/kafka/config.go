package kafka

import "orders-service/internal/config"

type Config struct {
	Brokers                     []string
	ConsumerGroup               string
	OrderCommandEventsTopic     string
	OrderCreatedEventType       string
	PaymentEventsTopic          string
	PaymentProcessedEventType   string
	OrderNotificationEventTopic string
	OrderStatusChangedEventType string
}

func NewKafkaConfig(conf *config.BrokerConfig) *Config {
	return &Config{
		Brokers:                     conf.Brokers,
		ConsumerGroup:               conf.ConsumerGroup,
		OrderCommandEventsTopic:     conf.OrderCommandEventsTopic,
		OrderCreatedEventType:       conf.OrderCreatedEventType,
		PaymentEventsTopic:          conf.PaymentEventsTopic,
		PaymentProcessedEventType:   conf.PaymentProcessedEventType,
		OrderNotificationEventTopic: conf.OrderNotificationEventTopic,
		OrderStatusChangedEventType: conf.OrderStatusChangedEventType,
	}
}
