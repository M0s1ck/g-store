package config

type BrokerConfig struct {
	Brokers               []string
	OrderEventsTopic      string
	OrderCreatedEventType string
	ConsumerGroup         string
}
