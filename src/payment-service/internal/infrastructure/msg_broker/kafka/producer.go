package kafka

import (
	"context"

	"github.com/segmentio/kafka-go"

	"payment-service/internal/domain/messages"
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
		Key:   msg.Id[:],
		Value: msg.Payload,
		Headers: []kafka.Header{
			{
				Key:   "message_id",
				Value: []byte(msg.Id.String()),
			},
			{
				Key:   "event_type",
				Value: []byte(msg.EventType),
			},
		},
	})
}
