package app

import (
	"log"
	"net/http"

	"orders-service/internal/config"
	mydelivery "orders-service/internal/delivery/http"
	"orders-service/internal/delivery/http/handlers"
	"orders-service/internal/infrastructure/db/postgres"
	"orders-service/internal/infrastructure/db/postgres/repository"
	servlogger "orders-service/internal/infrastructure/services/logger"
	"orders-service/internal/infrastructure/services/outbox"
	"orders-service/internal/usecase/create_order"
	"orders-service/internal/usecase/get_orders"
)

func Build(conf *config.Config) *http.Handler {
	psgConf := postgres.NewConfig(conf)
	logger := servlogger.NewSlogLogger()
	ordersDb, err := postgres.New(psgConf, logger)

	if err != nil {
		log.Fatal(err)
	}

	orderRepo := repository.NewOrderRepository(ordersDb)
	outboxRepo := repository.NewOutboxRepository(ordersDb)
	txManager := postgres.NewTxManager(ordersDb)

	outboxModelFactory := outbox.NewOutboxModelProtoFactory()

	getByIdUC := get_orders.NewGetByIdUsecase(orderRepo)
	getByUserUC := get_orders.NewGetByUserUsecase(orderRepo)

	createOrderUc := create_order.NewCreateOrderUsecase(txManager, orderRepo, outboxRepo, outboxModelFactory)

	orderHandler := handlers.NewOrderHandler(getByIdUC, getByUserUC, createOrderUc)

	router := mydelivery.NewRouter(&mydelivery.RouterDeps{
		OrderHandler: orderHandler,
	})

	return &router
}
