package kafka

import (
	"payment-service/internal/config"
)

type Config struct {
	Brokers                   []string
	ConsumerGroup             string
	OrderEventsTopic          string
	OrderCreatedEventType     string
	OrderCancelledEventType   string
	PaymentEventsTopic        string
	PaymentProcessedEventType string
	AllowedEventTypes         map[string]struct{}
}

func NewKafkaConfig(conf *config.BrokerConfig) *Config {
	return &Config{
		Brokers:                   conf.Brokers,
		ConsumerGroup:             conf.ConsumerGroup,
		OrderEventsTopic:          conf.OrderEventsTopic,
		OrderCreatedEventType:     conf.OrderCreatedEventType,
		OrderCancelledEventType:   conf.OrderCancelledEventType,
		PaymentEventsTopic:        conf.PaymentEventsTopic,
		PaymentProcessedEventType: conf.PaymentProcessedEventType,

		AllowedEventTypes: map[string]struct{}{
			conf.OrderCreatedEventType:   {},
			conf.OrderCancelledEventType: {},
		},
	}
}
