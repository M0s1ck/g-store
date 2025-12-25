package inbox

import (
	"context"
	"log"

	"payment-service/internal/usecase/common"
)

// MessageProcessor takes unprocessed messages and calls respective message handlers
type MessageProcessor struct {
	inboxRepo ProcessorRepo
	handlers  map[string]MessageHandler
	txManager common.TxManager
	batch     int
}

func NewMessageProcessor(
	repo ProcessorRepo,
	handlers []MessageHandler,
	txManager common.TxManager,
	batch int,
) *MessageProcessor {

	h := make(map[string]MessageHandler)
	for _, handler := range handlers {
		h[handler.EventType()] = handler
	}

	return &MessageProcessor{
		inboxRepo: repo,
		handlers:  h,
		txManager: txManager,
		batch:     batch,
	}
}

// ProcessBatch takes unprocessed messages and calls respective message handlers
func (p *MessageProcessor) ProcessBatch(ctx context.Context) error {
	return p.txManager.WithinTx(ctx, func(ctx context.Context) error {

		msgs, err := p.inboxRepo.GetUnprocessed(ctx, p.batch)
		if err != nil {
			return err
		}

		for _, msg := range msgs {
			handler, ok := p.handlers[msg.EventType]
			if !ok {
				log.Printf("no handler for event type %s", msg.EventType)
				continue
			}

			if err := handler.Handle(ctx, msg); err != nil {
				log.Printf("Error when handling message topic=%s id=%s",
					msg.Topic, msg.Id)
				return err
			}

			if err := p.inboxRepo.MarkProcessed(ctx, msg.Id); err != nil {
				return err
			}
		}

		return nil
	})
}
