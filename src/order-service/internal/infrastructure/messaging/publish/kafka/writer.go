package kafka

import (
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	batchSize    = 100
	batchTimeout = 500 * time.Millisecond
)

func NewKafkaWriter(conf *Config) *kafka.Writer {
	writer := &kafka.Writer{
		Addr:                   kafka.TCP(conf.Brokers...),
		Balancer:               &kafka.Hash{},
		Async:                  false,
		BatchTimeout:           batchTimeout,
		BatchSize:              batchSize,
		AllowAutoTopicCreation: false, // или true, если нужно
	}

	return writer
}
