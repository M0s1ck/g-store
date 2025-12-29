package app

import (
	"log"
	"net/http"
	"time"

	"orders-service/internal/config"
	"orders-service/internal/delivery/http"
	"orders-service/internal/delivery/http/handlers"
	"orders-service/internal/infrastructure/db/postgres"
	"orders-service/internal/infrastructure/db/postgres/repository"
	"orders-service/internal/infrastructure/msg_broker/kafka"
	"orders-service/internal/infrastructure/msg_broker/proto/mappers"
	"orders-service/internal/infrastructure/services/logger"
	"orders-service/internal/infrastructure/workers"
	"orders-service/internal/usecase/cancel_order"
	"orders-service/internal/usecase/common/outbox"
	"orders-service/internal/usecase/create_order"
	"orders-service/internal/usecase/event_handlers"
	"orders-service/internal/usecase/event_handlers/payment_processed"
	"orders-service/internal/usecase/get_orders"
	"orders-service/internal/usecase/order_update_status"
	"orders-service/internal/usecase/publish/outbox_publish"
)

func Build(conf *config.Config) (http.Handler, []workers.BackgroundWorker) {
	psgConf := postgres.NewConfig(conf)
	logger := services_logger.NewSlogLogger()
	ordersDb, err := postgres.New(psgConf, logger)

	if err != nil {
		log.Fatal(err)
	}

	kafkaConfig := msg_kafka.NewKafkaConfig(&conf.Broker)
	orderWriter := msg_kafka.NewKafkaWriter(kafkaConfig, kafkaConfig.OrderEventsTopic)
	orderProducer := msg_kafka.NewProducer(orderWriter, kafkaConfig.OrderCreatedEventType, kafkaConfig.OrderCancelledEventType, kafkaConfig.OrderStatusChangedEventType)
	paymentReader := msg_kafka.NewKafkaReader(kafkaConfig, kafkaConfig.PaymentEventsTopic)

	orderRepo := repository.NewOrderRepository(ordersDb)
	outboxRepo := repository.NewOutboxRepository(ordersDb)
	txManager := postgres.NewTxManager(ordersDb)

	getByIdUC := get_orders.NewGetByIdUsecase(orderRepo)
	getByUserUC := get_orders.NewGetByUserUsecase(orderRepo)

	ordCrMapper := proto_mappers.NewOrderCreatedPayloadMapper()
	payProcessedMapper := proto_mappers.NewPaymentProcessedPayloadMapper()
	ordCancelMapper := proto_mappers.NewOrderCancelledPayloadMapper()
	ordStChangedMapper := proto_mappers.NewOrderStatusChangedPayloadMapper()

	outboxMsgFactory := common_outbox.NewMessageFactory(
		ordCrMapper, ordCancelMapper, ordStChangedMapper,
		kafkaConfig.OrderCreatedEventType,
		kafkaConfig.OrderCancelledEventType,
		kafkaConfig.OrderStatusChangedEventType)

	updStatusPolicy := order_update_status.NewUpdateStatusPolicy()
	cancelPolicy := cancel_order.NewCancelPolicy()

	createOrderUc := create_order.NewCreateOrderUsecase(txManager, orderRepo, outboxRepo, outboxMsgFactory)
	cancelUC := cancel_order.NewCancelOrderUsecase(orderRepo, outboxRepo, txManager, outboxMsgFactory, cancelPolicy)
	updateStatusUC := order_update_status.NewUpdateOrderStatusUsecase(txManager, orderRepo, outboxRepo, updStatusPolicy, outboxMsgFactory)

	publishUC := outbox_publish.NewOutboxPublishUsecase(outboxRepo, orderProducer)

	orderHandler := handlers.NewOrderHandler(getByIdUC, getByUserUC, createOrderUc, cancelUC)
	staffOrderHandler := handlers.NewStaffOrderHandler(cancelUC, updateStatusUC)

	router := delivery_http.NewRouter(&delivery_http.RouterDeps{
		OrderHandler:      orderHandler,
		StaffOrderHandler: staffOrderHandler,
		Secrets:           conf.Secrets,
	})

	payProcessedEventHandler := payment_processed.NewPaymentProcessedEventHandler(
		cancelUC,
		updateStatusUC,
		payProcessedMapper,
		kafkaConfig.PaymentProcessedEventType)

	hands := []event_handlers.EventHandler{
		payProcessedEventHandler,
	}

	msgProcessor := event_handlers.NewEventMsgProcessor(hands)

	publishWorker := workers.NewOutboxPublishWorker(publishUC, 1*time.Second)
	kafkaConsumerWorker := workers.NewKafkaConsumerWorker(paymentReader, msgProcessor)

	bWorkers := []workers.BackgroundWorker{
		publishWorker,
		kafkaConsumerWorker,
	}

	return router, bWorkers
}
