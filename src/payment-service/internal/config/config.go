package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HTTP       HTTPConfig
	PaymentsDB DBConfig
	Broker     BrokerConfig
}

func Load() (*Config, error) {
	cfg := &Config{
		HTTP: HTTPConfig{
			Addr: getEnv("HTTP_ADDR", ":8080"),
		},
		PaymentsDB: DBConfig{
			Host:     os.Getenv("DB_HOST"),
			Port:     os.Getenv("DB_PORT"),
			Name:     os.Getenv("DB_NAME"),
			User:     os.Getenv("DB_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),

			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", 10),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getEnvDuration("DB_CONN_MAX_LIFETIME", 30*time.Second),
		},
		Broker: BrokerConfig{
			Brokers:                   strings.Split(os.Getenv("BROKER_HOST"), ","),
			ConsumerGroup:             os.Getenv("BROKER_CONSUMER_GROUP"),
			OrderEventsTopic:          os.Getenv("BROKER_ORDER_EVENTS_TOPIC"),
			OrderCreatedEventType:     os.Getenv("BROKER_ORDER_CREATED_EVENT_TYPE"),
			PaymentEventsTopic:        os.Getenv("BROKER_PAYMENT_EVENTS_TOPIC"),
			PaymentProcessedEventType: os.Getenv("BROKER_PAYMENT_PROCESSED_EVENT_TYPE"),
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

func getEnvInt(key string, def int) int {
	if str := os.Getenv(key); str != "" {
		num, err := strconv.Atoi(str)
		if err != nil {
			return def
		}
		return num
	}
	return def
}

func getEnvDuration(key string, def time.Duration) time.Duration {
	if str := os.Getenv(key); str != "" {
		num, err := strconv.Atoi(str)
		if err != nil {
			return def
		}
		return time.Duration(num) * time.Second
	}

	return def
}
