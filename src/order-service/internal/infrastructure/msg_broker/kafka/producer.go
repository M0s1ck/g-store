package msg_kafka

import (
	"context"

	"github.com/segmentio/kafka-go"

	"orders-service/internal/domain/messages"
)

type Producer struct {
	ordEventWriter      *kafka.Writer
	ordCreatedEvtType   string
	ordCancelledEvtType string
	ordStChangedEvtType string
}

func NewProducer(
	orderWriter *kafka.Writer,
	ordCreatedEvtType string,
	ordCancelledEvtType string,
	ordStChangedType string,
) *Producer {

	return &Producer{
		ordEventWriter:      orderWriter,
		ordCreatedEvtType:   ordCreatedEvtType,
		ordCancelledEvtType: ordCancelledEvtType,
		ordStChangedEvtType: ordStChangedType,
	}
}

func (p *Producer) Publish(ctx context.Context, msg *messages.OutboxMessage) error {
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

	return p.ordEventWriter.WriteMessages(ctx, kafkaMsg)
}
