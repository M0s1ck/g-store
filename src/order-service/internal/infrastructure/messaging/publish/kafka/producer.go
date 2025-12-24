package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"

	"orders-service/internal/domain/messages"
)

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(writer *kafka.Writer) *Producer {
	return &Producer{
		writer: writer,
	}
}

func (k *Producer) Publish(ctx context.Context, msg *messages.OutboxMessage) error {
	return k.writer.WriteMessages(ctx, kafka.Message{
		Topic: msg.EventType,
		Key:   msg.Id[:],
		Value: msg.Payload,
	})
}
