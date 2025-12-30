package background_workers

import (
	"context"
	"log"
	"time"

	"payment-service/internal/usecase/inbox"
)

// InboxProcessWorker calls message processor in a loop
type InboxProcessWorker struct {
	processor *inbox.MessageProcessor
	interval  time.Duration
}

func NewInboxProcessWorker(
	processor *inbox.MessageProcessor,
	interval time.Duration,
) *InboxProcessWorker {
	return &InboxProcessWorker{
		processor: processor,
		interval:  interval,
	}
}

func (w *InboxProcessWorker) Run(ctx context.Context) error {
	log.Printf("Starting inbox processor worker")
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return nil

		case <-ticker.C:
			if err := w.processor.ProcessBatch(ctx); err != nil {
				log.Printf("inbox worker error: %v", err)
			}
		}
	}
}
