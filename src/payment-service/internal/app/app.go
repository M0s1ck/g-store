package app

import (
	"log"
	"net/http"
	"time"

	"payment-service/internal/config"
	httpdelivery "payment-service/internal/delivery/http"
	"payment-service/internal/delivery/http/handlers"
	"payment-service/internal/infrastructure/background_workers"
	"payment-service/internal/infrastructure/db/postgres"
	"payment-service/internal/infrastructure/db/postgres/repository"
	kafka2 "payment-service/internal/infrastructure/msg_broker/kafka"
	"payment-service/internal/infrastructure/msg_broker/proto/mappers"
	servlogger "payment-service/internal/infrastructure/services/logger"
	"payment-service/internal/usecase/common/outbox"
	"payment-service/internal/usecase/create_account"
	"payment-service/internal/usecase/event_handlers/order_cancelled"
	"payment-service/internal/usecase/event_handlers/order_created"
	"payment-service/internal/usecase/get_account"
	"payment-service/internal/usecase/inbox"
	outboxuc "payment-service/internal/usecase/outbox"
	"payment-service/internal/usecase/pay_for_order"
	"payment-service/internal/usecase/refund_order"
	"payment-service/internal/usecase/top_up"
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
	accountRepo := repository.NewAccountRepository(paymentsDb)
	transactionRepo := repository.NewBalanceTransactionRepository(paymentsDb)
	outboxRepo := repository.NewOutboxRepository(paymentsDb)

	getAccUC := get_account.NewGetByIdUsecase(accountRepo)
	createAccUC := create_account.NewCreateAccountUsecase(accountRepo)
	topUpUC := top_up.NewTopUpUsecase(accountRepo, transactionRepo, txManager)
	refundUC := refund_order.NewRefundUsecase(accountRepo, transactionRepo, txManager)

	accountHandler := handlers.NewAccountHandler(getAccUC, createAccUC, topUpUC)

	router := httpdelivery.NewRouter(&httpdelivery.RouterDeps{
		AccountHandler: accountHandler,
	})

	ordCreatedMapper := proto_mappers.NewOrderCreatedPayloadMapper()
	paymentProcMapper := proto_mappers.NewPaymentProcessedPayloadMapper()
	ordCancelMapper := proto_mappers.NewOrderCancelledPayloadMapper()

	kafkaConfig := kafka2.NewKafkaConfig(&conf.Broker)
	kafkaReader := kafka2.NewKafkaOrderCreatedReader(kafkaConfig)
	kafkaPaymentsWriter := kafka2.NewKafkaWriter(kafkaConfig, kafkaConfig.PaymentEventsTopic)
	kafkaProducer := kafka2.NewProducer(kafkaPaymentsWriter)

	kafkaConsumerWorker := kafka2.NewInboxKafkaConsumerWorker(inboxRepo, kafkaReader, kafkaConfig.AllowedEventTypes)

	outboxMsgFactory := outbox.NewOutboxMessageFactory(paymentProcMapper, kafkaConfig.PaymentProcessedEventType)

	payForOrderUC := pay_for_order.NewPayUsecase(accountRepo, transactionRepo, txManager, outboxRepo, outboxMsgFactory)

	ordCreatedEvtHandler := order_created.NewEventHandler(payForOrderUC, ordCreatedMapper, kafkaConfig.OrderCreatedEventType)
	ordCancelledEvtHandler := order_cancelled.NewEventHandler(refundUC, ordCancelMapper, kafkaConfig.OrderCancelledEventType)

	messageHandlers := []inbox.MessageHandler{
		ordCreatedEvtHandler,
		ordCancelledEvtHandler,
	}

	processor := inbox.NewMessageProcessor(inboxRepo, messageHandlers, txManager, 10)

	inboxProcessWorker := background_workers.NewInboxProcessWorker(processor, 1*time.Second)

	publishUC := outboxuc.NewPublishUsecase(outboxRepo, kafkaProducer)
	outboxPublishWorker := background_workers.NewOutboxPublishWorker(publishUC, 1*time.Second)

	backgroundWorkers := []background_workers.BackgroundWorker{
		kafkaConsumerWorker,
		inboxProcessWorker,
		outboxPublishWorker,
	}

	return router, backgroundWorkers
}
