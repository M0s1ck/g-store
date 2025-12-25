package app

import (
	"log"
	"net/http"
	"time"

	"orders-service/internal/config"
	mydelivery "orders-service/internal/delivery/http"
	"orders-service/internal/delivery/http/handlers"
	"orders-service/internal/infrastructure/db/postgres"
	"orders-service/internal/infrastructure/db/postgres/repository"
	"orders-service/internal/infrastructure/messaging/publish/kafka"
	servlogger "orders-service/internal/infrastructure/services/logger"
	"orders-service/internal/infrastructure/services/outbox"
	"orders-service/internal/infrastructure/workers"
	"orders-service/internal/usecase/create_order"
	"orders-service/internal/usecase/get_orders"
	"orders-service/internal/usecase/publish/outbox_publish"
)

func Build(conf *config.Config) (*http.Handler, *workers.OutboxPublishWorker) {
	psgConf := postgres.NewConfig(conf)
	logger := servlogger.NewSlogLogger()
	ordersDb, err := postgres.New(psgConf, logger)

	if err != nil {
		log.Fatal(err)
	}

	kafkaConfig := kafka.NewKafkaConfig(&conf.Broker)
	orderWriter := kafka.NewKafkaWriter(kafkaConfig, kafkaConfig.OrderEventsTopic)
	kafkaProducer := kafka.NewProducer(orderWriter)

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

	return &router, publishWorker
}
