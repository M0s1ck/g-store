package kafka

import "github.com/segmentio/kafka-go"

const commitInterval = 0 // for manual commits

func NewKafkaOrderCreatedReader(conf *Config) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        conf.Brokers,
		Topic:          conf.OrderEventsTopic,
		GroupID:        conf.ConsumerGroup,
		CommitInterval: commitInterval,
	})
}
