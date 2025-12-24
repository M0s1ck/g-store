package kafka

import "orders-service/internal/config"

type Config struct {
	Brokers           []string
	OrderCreatedTopic string
}

func NewKafkaConfig(conf *config.BrokerConfig) *Config {
	return &Config{
		Brokers:           conf.Brokers,
		OrderCreatedTopic: conf.OrderCreatedTopic,
	}
}
