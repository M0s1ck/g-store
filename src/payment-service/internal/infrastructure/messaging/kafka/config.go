package kafka

import (
	"payment-service/internal/config"
)

type Config struct {
	Brokers                   []string
	ConsumerGroup             string
	OrderEventsTopic          string
	OrderCreatedEventType     string
	PaymentEventsTopic        string
	PaymentProcessedEventType string
}

func NewKafkaConfig(conf *config.BrokerConfig) *Config {
	return &Config{
		Brokers:                   conf.Brokers,
		ConsumerGroup:             conf.ConsumerGroup,
		OrderEventsTopic:          conf.OrderEventsTopic,
		OrderCreatedEventType:     conf.OrderCreatedEventType,
		PaymentEventsTopic:        conf.PaymentEventsTopic,
		PaymentProcessedEventType: conf.PaymentProcessedEventType,
	}
}
