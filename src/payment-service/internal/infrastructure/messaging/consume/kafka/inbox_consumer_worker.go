package kafka

import (
	"context"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/segmentio/kafka-go"

	"payment-service/internal/domain/messages"
	"payment-service/internal/usecase/inbox"
)

const (
	messageIdHeaderName = "message_id"
	eventTypeHeaderName = "event_type"
)

// InboxKafkaConsumerWorker background worker to fetch messages from kafka and write them to inbox
type InboxKafkaConsumerWorker struct {
	repo   inbox.BrokerConsumerRepo
	reader *kafka.Reader
}

func NewInboxKafkaConsumerWorker(repo inbox.BrokerConsumerRepo, reader *kafka.Reader) *InboxKafkaConsumerWorker {
	return &InboxKafkaConsumerWorker{repo: repo, reader: reader}
}

func (c *InboxKafkaConsumerWorker) Run(ctx context.Context) error {
	for {
		msg, err := c.reader.FetchMessage(ctx)
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

		inboxMsg := messages.InboxMessage{
			Id:        mesId,
			Topic:     msg.Topic,
			EventType: eventType,
			Key:       msg.Key,
			Payload:   msg.Value,
		}

		if err := c.repo.SaveIdempotent(ctx, inboxMsg); err != nil {
			log.Printf("failed to save inbox msg: id=%v event=%v: %v",
				inboxMsg.Id, inboxMsg.EventType, err)
			continue
		}

		err = c.reader.CommitMessages(ctx, msg)
		if err != nil {
			log.Printf("failed to commit msg: id=%v event=%v: %v",
				inboxMsg.Id, inboxMsg.EventType, err)
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
