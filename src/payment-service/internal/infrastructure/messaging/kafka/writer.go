package kafka

import (
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	batchSize    = 100
	batchTimeout = 500 * time.Millisecond
)

func NewKafkaWriter(conf *Config, topic string) *kafka.Writer {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(conf.Brokers...),
		Balancer:               &kafka.Hash{},
		Topic:                  topic,
		Async:                  false,
		BatchTimeout:           batchTimeout,
		BatchSize:              batchSize,
		AllowAutoTopicCreation: false,
	}

	return writer
}
