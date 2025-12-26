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
	kafka2 "payment-service/internal/infrastructure/messaging/kafka"
	servlogger "payment-service/internal/infrastructure/services/logger"
	"payment-service/internal/infrastructure/services/proto/proto_order_created"
	"payment-service/internal/infrastructure/services/proto/proto_payment_processed"
	"payment-service/internal/usecase/common/outbox"
	"payment-service/internal/usecase/create_account"
	"payment-service/internal/usecase/get_account"
	"payment-service/internal/usecase/inbox"
	"payment-service/internal/usecase/order_created"
	outboxuc "payment-service/internal/usecase/outbox"
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
	balanceTransactionRepo := repository.NewBalanceTransactionRepository(paymentsDb)
	outboxRepo := repository.NewOutboxRepository(paymentsDb)

	getAccUC := get_account.NewGetByIdUsecase(accountRepo)
	createAccUC := create_account.NewCreateAccountUsecase(accountRepo)
	topUpUC := top_up.NewTopUpUsecase(accountRepo, balanceTransactionRepo, txManager)

	accountHandler := handlers.NewAccountHandler(getAccUC, createAccUC, topUpUC)

	router := httpdelivery.NewRouter(&httpdelivery.RouterDeps{
		AccountHandler: accountHandler,
	})

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
