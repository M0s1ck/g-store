package workers

import (
	"context"
	"log"
	"time"

	"orders-service/internal/usecase/publish/outbox_publish"
)

const batchSize = 100

type OutboxPublishWorker struct {
	publisher *outbox_publish.OutboxPublishUsecase
	ticker    *time.Ticker
}

func NewOutboxPublishWorker(
	publisher *outbox_publish.OutboxPublishUsecase,
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
