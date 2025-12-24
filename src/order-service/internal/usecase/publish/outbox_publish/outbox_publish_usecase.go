package outbox_publish

import (
	"context"
	"log"

	"orders-service/internal/domain/messages"
)

type OutboxPublishUsecase struct {
	repo     OutboxRepoPublisher
	producer EventProducer
}

func NewOutboxPublishUsecase(
	repo OutboxRepoPublisher,
	producer EventProducer,
) *OutboxPublishUsecase {

	return &OutboxPublishUsecase{
		repo:     repo,
		producer: producer,
	}
}

func (p *OutboxPublishUsecase) Publish(ctx context.Context, batchSize int) error {
	msgs, err := p.repo.GetUnsent(ctx, batchSize)
	if err != nil {
		return err
	}

	for _, msg := range msgs {
		err := p.producer.Publish(ctx, &msg)

		if err != nil {
			p.failedToPublish(ctx, &msg, err)
		}

		err = p.repo.MarkAsSent(ctx, msg.Id)
		if err != nil {
			log.Printf("failed to mark outbox message as sent: %v", err)
		}
	}

	return nil
}

func (p *OutboxPublishUsecase) failedToPublish(
	ctx context.Context,
	msg *messages.OutboxMessage,
	err error,
) {

	incErr := p.repo.IncrementRetry(ctx, msg.AggregateID)
	if incErr != nil {
		log.Printf("Error incrementing retry: %v", incErr)
	}

	log.Printf("Error publishing: %v : %v : %v: %v",
		msg.EventType, msg.Aggregate, msg.AggregateID, err.Error())
}
