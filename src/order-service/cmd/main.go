package main

import (
	"context"
	"log"
	"net/http"

	"orders-service/internal/app"
	"orders-service/internal/config"
)

// @title Order Service API
// @version 1.0
// @description Order microservice
// @BasePath /v1
// @schemes http https
func main() {
	log.Println("Service is starting...")

	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	handler, publishWorker := app.Build(conf)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go publishWorker.Run(ctx)

	log.Println("Server started!")

	err = http.ListenAndServe(":8080", *handler)
	if err != nil {
		panic(err)
	}
}
