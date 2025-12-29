package app

import (
	"log"
	"net/http"
	"orders-service/internal/usecase/cancel_order"
	"time"

	"orders-service/internal/config"
	mydelivery "orders-service/internal/delivery/http"
	"orders-service/internal/delivery/http/handlers"
	"orders-service/internal/delivery/proto/order_created"
	protostatuschanged "orders-service/internal/delivery/proto/order_status_changed"
	protopayment "orders-service/internal/delivery/proto/payment_processed"
	"orders-service/internal/infrastructure/db/postgres"
	"orders-service/internal/infrastructure/db/postgres/repository"
	kafka2 "orders-service/internal/infrastructure/messaging/kafka"
	servlogger "orders-service/internal/infrastructure/services/logger"
	"orders-service/internal/infrastructure/workers"
	comoutbox "orders-service/internal/usecase/common/outbox"
	"orders-service/internal/usecase/create_order"
	"orders-service/internal/usecase/event_handlers"
	"orders-service/internal/usecase/event_handlers/payment_processed"
	"orders-service/internal/usecase/get_orders"
	"orders-service/internal/usecase/order_update_status"
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
	orderNotifWriter := kafka2.NewKafkaWriter(kafkaConfig, kafkaConfig.OrderNotificationEventTopic)
	orderProducer := kafka2.NewProducer(orderWriter, orderNotifWriter, kafkaConfig.OrderCreatedEventType, kafkaConfig.OrderStatusChangedEventType)
	paymentReader := kafka2.NewKafkaReader(kafkaConfig, kafkaConfig.PaymentEventsTopic)

	orderRepo := repository.NewOrderRepository(ordersDb)
	outboxRepo := repository.NewOutboxRepository(ordersDb)
	txManager := postgres.NewTxManager(ordersDb)

	getByIdUC := get_orders.NewGetByIdUsecase(orderRepo)
	getByUserUC := get_orders.NewGetByUserUsecase(orderRepo)

	ordCrMapper := order_created.NewPayloadMapper()
	payProcessedMapper := protopayment.NewPayloadMapper()
	ordStChangedMapper := protostatuschanged.NewPayloadMapper()

	outboxMsgFactory := comoutbox.NewMessageFactory(
		ordCrMapper, ordStChangedMapper,
		kafkaConfig.OrderCreatedEventType, kafkaConfig.OrderStatusChangedEventType)

	createOrderUc := create_order.NewCreateOrderUsecase(txManager, orderRepo, outboxRepo, outboxMsgFactory)

	updStatusPolicy := order_update_status.NewUpdateStatusPolicy()
	updateStatusUC := order_update_status.NewUpdateOrderStatusUsecase(txManager, orderRepo, outboxRepo, updStatusPolicy, outboxMsgFactory)

	cancelPolicy := cancel_order.NewCancelPolicy()
	cancelUC := cancel_order.NewCancelOrderUsecase(orderRepo, cancelPolicy)

	publishUC := outbox_publish.NewOutboxPublishUsecase(outboxRepo, orderProducer)

	orderHandler := handlers.NewOrderHandler(getByIdUC, getByUserUC, createOrderUc, cancelUC)
	staffOrderHandler := handlers.NewStaffOrderHandler(cancelUC, updateStatusUC)

	router := mydelivery.NewRouter(&mydelivery.RouterDeps{
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
