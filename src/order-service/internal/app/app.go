package app

import (
	"log"
	"net/http"
	"time"

	"orders-service/internal/config"
	mydelivery "orders-service/internal/delivery/http"
	"orders-service/internal/delivery/http/handlers"
	"orders-service/internal/delivery/proto/outbox"
	protopayment "orders-service/internal/delivery/proto/payment_processed"
	"orders-service/internal/infrastructure/db/postgres"
	"orders-service/internal/infrastructure/db/postgres/repository"
	kafka2 "orders-service/internal/infrastructure/messaging/kafka"
	servlogger "orders-service/internal/infrastructure/services/logger"
	"orders-service/internal/infrastructure/workers"
	"orders-service/internal/usecase/create_order"
	"orders-service/internal/usecase/event_handlers"
	"orders-service/internal/usecase/event_handlers/payment_processed"
	"orders-service/internal/usecase/get_orders"
	"orders-service/internal/usecase/publish/outbox_publish"
)

func Build(conf *config.Config) (http.Handler, []workers.BackgroundWorker) {
	psgConf := postgres.NewConfig(conf)
	logger := servlogger.NewSlogLogger()
	ordersDb, err := postgres.New(psgConf, logger)

	if err != nil {
		log.Fatal(err)
	}

	kafkaConfig := kafka2.NewKafkaConfig(&conf.Broker)
	orderWriter := kafka2.NewKafkaWriter(kafkaConfig, kafkaConfig.OrderCommandEventsTopic)
	kafkaProducer := kafka2.NewProducer(orderWriter)

	orderRepo := repository.NewOrderRepository(ordersDb)
	outboxRepo := repository.NewOutboxRepository(ordersDb)
	txManager := postgres.NewTxManager(ordersDb)

	outboxModelFactory := outbox.NewOutboxModelProtoFactory(kafkaConfig.OrderCreatedEventType)

	getByIdUC := get_orders.NewGetByIdUsecase(orderRepo)
	getByUserUC := get_orders.NewGetByUserUsecase(orderRepo)

	createOrderUc := create_order.NewCreateOrderUsecase(txManager, orderRepo, outboxRepo, outboxModelFactory)

	publishUC := outbox_publish.NewOutboxPublishUsecase(outboxRepo, kafkaProducer)

	publishWorker := workers.NewOutboxPublishWorker(publishUC, 1*time.Second)

	orderHandler := handlers.NewOrderHandler(getByIdUC, getByUserUC, createOrderUc)

	router := mydelivery.NewRouter(&mydelivery.RouterDeps{
		OrderHandler: orderHandler,
	})

	paymentReader := kafka2.NewKafkaReader(kafkaConfig, kafkaConfig.PaymentEventsTopic)

	protoPaymentProcessedMapper := protopayment.NewPayloadMapper()

	paymentProcessedEventHandler := payment_processed.NewPaymentProcessedEventHandler(
		kafkaConfig.PaymentProcessedEventType, orderRepo, protoPaymentProcessedMapper)

	hands := []event_handlers.EventHandler{
		paymentProcessedEventHandler,
	}

	msgProcessor := event_handlers.NewEventMsgProcessor(hands)

	kafkaConsumerWorker := workers.NewKafkaConsumerWorker(paymentReader, msgProcessor)

	bWorkers := []workers.BackgroundWorker{
		publishWorker,
		kafkaConsumerWorker,
	}

	return router, bWorkers
}
