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

	getByIdUC := get_orders.NewGetByIdUsecase(orderRepo)
	getByUserUC := get_orders.NewGetByUserUsecase(orderRepo)

	orderHandler := handlers.NewOrderHandler(getByIdUC, getByUserUC)

	router := mydelivery.NewRouter(&mydelivery.RouterDeps{
		OrderHandler: orderHandler,
	})

	return &router
}
