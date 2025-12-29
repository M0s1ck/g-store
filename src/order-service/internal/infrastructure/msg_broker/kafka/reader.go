package msg_kafka

import "github.com/segmentio/kafka-go"

const commitInterval = 0 // for manual commits

func NewKafkaReader(conf *Config, topic string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:        conf.Brokers,
		Topic:          topic,
		GroupID:        conf.ConsumerGroup,
		CommitInterval: commitInterval,
	})
}
