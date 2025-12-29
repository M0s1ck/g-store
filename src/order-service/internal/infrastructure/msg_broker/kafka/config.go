package msg_kafka

import "orders-service/internal/config"

type Config struct {
	Brokers                     []string
	ConsumerGroup               string
	OrderEventsTopic            string
	OrderCreatedEventType       string
	OrderStatusChangedEventType string
	OrderCancelledEventType     string
	PaymentEventsTopic          string
	PaymentProcessedEventType   string
}

func NewKafkaConfig(conf *config.BrokerConfig) *Config {
	return &Config{
		Brokers:                     conf.Brokers,
		ConsumerGroup:               conf.ConsumerGroup,
		OrderEventsTopic:            conf.OrderEventsTopic,
		OrderCreatedEventType:       conf.OrderCreatedEventType,
		PaymentEventsTopic:          conf.PaymentEventsTopic,
		PaymentProcessedEventType:   conf.PaymentProcessedEventType,
		OrderCancelledEventType:     conf.OrderCancelledEventType,
		OrderStatusChangedEventType: conf.OrderStatusChangedEventType,
	}
}
