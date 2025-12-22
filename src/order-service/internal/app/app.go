package app

import (
	"net/http"

	"orders-service/internal/config"
	mydelivery "orders-service/internal/delivery/http"
	"orders-service/internal/delivery/http/handlers"
)

func Build(conf *config.Config) *http.Handler {
	orderHandler := handlers.NewOrderHandler()

	router := mydelivery.NewRouter(&mydelivery.RouterDeps{
		OrderHandler: orderHandler,
	})

	return &router
}
