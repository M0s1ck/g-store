package app

import (
	"log"
	"net/http"
	kafka2 "payment-service/internal/infrastructure/messaging/kafka"
	"payment-service/internal/infrastructure/services/proto/proto_order_created"
	"payment-service/internal/infrastructure/services/proto/proto_payment_processed"
	"payment-service/internal/usecase/common/outbox"
	outboxuc "payment-service/internal/usecase/outbox"
	"time"

	"payment-service/internal/config"
	"payment-service/internal/delivery"
	"payment-service/internal/infrastructure/background_workers"
	"payment-service/internal/infrastructure/db/postgres"
	"payment-service/internal/infrastructure/db/postgres/repository"
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
	accountRepo := repository.NewAccountRepository(paymentsDb)
	balanceTransactionRepo := repository.NewBalanceTransactionRepository(paymentsDb)
	outboxRepo := repository.NewOutboxRepository(paymentsDb)

	router := delivery.NewRouter(&delivery.RouterDeps{})

	protoOrderMapper := proto_order_created.NewPayloadMapper()
	protoPaymentMapper := proto_payment_processed.NewPayloadMapper()

	kafkaConfig := kafka2.NewKafkaConfig(&conf.Broker)
	kafkaReader := kafka2.NewKafkaOrderCreatedReader(kafkaConfig)
	kafkaPaymentsWriter := kafka2.NewKafkaWriter(kafkaConfig, kafkaConfig.PaymentEventsTopic)
	kafkaProducer := kafka2.NewProducer(kafkaPaymentsWriter)

	kafkaConsumerWorker := kafka2.NewInboxKafkaConsumerWorker(inboxRepo, kafkaReader)

	outboxMsgFactory := outbox.NewOutboxMessageFactory(protoPaymentMapper, kafkaConfig.PaymentProcessedEventType)

	orderCreatedEventHandler := order_created.NewOrderCreatedEventHandler(
		accountRepo,
		balanceTransactionRepo,
		outboxRepo,
		outboxMsgFactory,
		protoOrderMapper,
		kafkaConfig.OrderCreatedEventType,
	)

	messageHandlers := []inbox.MessageHandler{
		orderCreatedEventHandler,
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
