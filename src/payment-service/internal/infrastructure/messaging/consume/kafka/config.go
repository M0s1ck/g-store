package kafka

import (
	"payment-service/internal/config"
)

type Config struct {
	Brokers               []string
	OrderEventsTopic      string
	OrderCreatedEventType string
	ConsumerGroup         string
}

func NewKafkaConfig(conf *config.BrokerConfig) *Config {
	return &Config{
		Brokers:               conf.Brokers,
		OrderEventsTopic:      conf.OrderEventsTopic,
		OrderCreatedEventType: conf.OrderCreatedEventType,
		ConsumerGroup:         conf.ConsumerGroup,
	}
}
