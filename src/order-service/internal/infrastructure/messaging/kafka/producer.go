package kafka

import (
	"context"
	"fmt"

	"github.com/segmentio/kafka-go"

	"orders-service/internal/domain/messages"
)

type Producer struct {
	ordCommandWriter      *kafka.Writer
	ordNotificationWriter *kafka.Writer
	ordCreatedEvtType     string
	ordStChangedEvtType   string
}

func NewProducer(
	orderWriter *kafka.Writer,
	orderNotificationWriter *kafka.Writer,
	ordCreatedEvtType string,
	ordStChangedType string,
) *Producer {

	return &Producer{
		ordCommandWriter:      orderWriter,
		ordNotificationWriter: orderNotificationWriter,
		ordCreatedEvtType:     ordCreatedEvtType,
		ordStChangedEvtType:   ordStChangedType,
	}
}

func (k *Producer) Publish(ctx context.Context, msg *messages.OutboxMessage) error {
	kafkaMsg := kafka.Message{
		Key:   []byte(msg.Key),
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
	}

	switch msg.EventType {
	case k.ordCreatedEvtType:
		return k.ordCommandWriter.WriteMessages(ctx, kafkaMsg)

	case k.ordStChangedEvtType:
		return k.ordNotificationWriter.WriteMessages(ctx, kafkaMsg)
	}

	return fmt.Errorf("unknown event type: %s", msg.EventType)
}
