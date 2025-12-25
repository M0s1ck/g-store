package background_workers

import (
	"context"
	"log"
	"payment-service/internal/usecase/outbox"
	"time"
)

const batchSize = 100

// OutboxPublishWorker calls publish usecase in a for loop
type OutboxPublishWorker struct {
	publisher *outbox.PublishUsecase
	ticker    *time.Ticker
}

func NewOutboxPublishWorker(
	publisher *outbox.PublishUsecase,
	interval time.Duration,
) *OutboxPublishWorker {
	return &OutboxPublishWorker{
		publisher: publisher,
		ticker:    time.NewTicker(interval),
	}
}

func (w *OutboxPublishWorker) Run(ctx context.Context) error {
	log.Println("Starting outbox publish worker")

	for {
		select {
		case <-ctx.Done():
			return nil
		case <-w.ticker.C:
			if err := w.publisher.Publish(ctx, batchSize); err != nil {
				log.Println("Failed to publish outbox worker:", err)
				return err
			}
		}
	}
}
