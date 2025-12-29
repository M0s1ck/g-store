package config

import (
	"os"
	"strings"
)

type Config struct {
	Broker BrokerConfig
	Net    NetConfig
}

func Load() (*Config, error) {
	cfg := &Config{
		Net: NetConfig{
			Addr: getEnv("NET_ADDR", ":8080"),
		},
		Broker: BrokerConfig{
			Brokers:                     strings.Split(os.Getenv("BROKER_HOST"), ","),
			ConsumerGroup:               os.Getenv("BROKER_CONSUMER_GROUP"),
			OrderNotificationEventTopic: os.Getenv("BROKER_ORDER_NOTIFICATION_EVENTS_TOPIC"),
			OrderStatusChangedEventType: os.Getenv("BROKER_ORDER_STATUS_CHANGED_EVENT_TYPE"),
		},
	}

	return cfg, nil
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
