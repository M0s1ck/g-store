package app

import (
	"log"
	"net/http"
	"time"

	"payment-service/internal/config"
	"payment-service/internal/delivery"
	"payment-service/internal/infrastructure/background_workers"
	"payment-service/internal/infrastructure/db/postgres"
	"payment-service/internal/infrastructure/db/postgres/repository"
	"payment-service/internal/infrastructure/messaging/consume/kafka"
	servlogger "payment-service/internal/infrastructure/services/logger"
	"payment-service/internal/usecase/inbox"
	"payment-service/internal/usecase/order_created"
)

func Build(conf *config.Config) (http.Handler, []background_workers.BackgroundWorker) {
	logger := servlogger.NewSlogLogger()
	psgConf := postgres.NewConfig(conf)
	paymentsDb, err := postgres.New(psgConf, logger)

	if err != nil {
		log.Fatal(err)
	}

	txManager := postgres.NewTxManager(paymentsDb)
	inboxRepo := repository.NewInboxRepository(paymentsDb)

	router := delivery.NewRouter(&delivery.RouterDeps{})

	kafkaConfig := kafka.NewKafkaConfig(&conf.Broker)
	kafkaReader := kafka.NewKafkaOrderCreatedReader(kafkaConfig)

	kafkaConsumerWorker := kafka.NewInboxKafkaConsumerWorker(inboxRepo, kafkaReader)

	orderCreatedEventHandler :=
		order_created.NewOrderCreatedEventHandler(kafkaConfig.OrderCreatedEventType)

	messageHandlers := []inbox.MessageHandler{
		orderCreatedEventHandler,
	}

	processor := inbox.NewMessageProcessor(inboxRepo, messageHandlers, txManager, 10)

	inboxProcessWorker := background_workers.NewInboxProcessWorker(
		processor,
		1*time.Second,
	)

	backgroundWorkers := []background_workers.BackgroundWorker{
		inboxProcessWorker,
		kafkaConsumerWorker,
	}

	return router, backgroundWorkers
}
