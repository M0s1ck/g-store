package workers

import (
	"context"
	"errors"
	"log"
	"orders-service/internal/usecase/event_handlers"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"
)

const (
	messageIdHeaderName = "message_id"
	eventTypeHeaderName = "event_type"
)

type KafkaConsumerWorker struct {
	reader       *kafka.Reader
	msgProcessor *event_handlers.EventMsgProcessor
}

func NewKafkaConsumerWorker(reader *kafka.Reader, msgProcessor *event_handlers.EventMsgProcessor) *KafkaConsumerWorker {
	return &KafkaConsumerWorker{
		reader:       reader,
		msgProcessor: msgProcessor,
	}
}

func (w *KafkaConsumerWorker) Run(ctx context.Context) error {
	for {
		msg, err := w.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) { // expected shutdown
				return nil
			}
			log.Printf("Error fetching message from Kafka: %v", err)
			continue
		}

		mesId, err := getUUID(msg, messageIdHeaderName)
		if err != nil {
			log.Printf("invalid msg id: %v", err)
			continue
		}

		eventType := getHeader(msg, eventTypeHeaderName)
		if eventType == "" {
			log.Printf("missing event_type for msg %v", mesId)
			continue
		}

		err = w.msgProcessor.Process(ctx, msg.Value, eventType)
		if err != nil {
			log.Printf("error while proccesing msg evt=%s : %s", eventType, err)
			continue
		}

		err = w.reader.CommitMessages(ctx, msg)
		if err != nil {
			log.Printf("error while committing msg evt=%s : %s", eventType, err)
			continue
		}
	}
}

func getHeader(msg kafka.Message, key string) string {
	for _, h := range msg.Headers {
		if h.Key == key {
			return string(h.Value)
		}
	}
	return ""
}

func getUUID(msg kafka.Message, key string) (uuid.UUID, error) {
	str := getHeader(msg, key)
	return uuid.Parse(str)
}
