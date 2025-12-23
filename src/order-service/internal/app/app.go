package app

import (
	"log"
	"net/http"
	"orders-service/internal/config"
	mydelivery "orders-service/internal/delivery/http"
	"orders-service/internal/delivery/http/handlers"
	"orders-service/internal/infrastructure/db/postgres"
	"orders-service/internal/infrastructure/db/postgres/repository"
	"orders-service/internal/usecase/get_orders"
)

func Build(conf *config.Config) *http.Handler {
	psgConf := postgres.NewConfig(conf)
	ordersDb, err := postgres.New(psgConf)
	if err != nil {
		log.Fatal(err)
	}

	orderRepo := repository.NewOrderRepository(ordersDb)

	getUC := get_orders.NewGetOrdersUsecase(orderRepo)

	orderHandler := handlers.NewOrderHandler(getUC)

	router := mydelivery.NewRouter(&mydelivery.RouterDeps{
		OrderHandler: orderHandler,
	})

	return &router
}
