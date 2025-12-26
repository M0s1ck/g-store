package main

import (
	"context"
	"log"
	"net/http"

	"orders-service/internal/app"
	"orders-service/internal/config"
	"orders-service/internal/infrastructure/workers"
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

	handler, backWorkers := app.Build(conf)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// launching background workers
	for _, worker := range backWorkers {
		go func(worker workers.BackgroundWorker) {
			err := worker.Run(ctx)
			if err != nil {
				log.Printf("worker stopped with error: %v", err)
				cancel() // if one has error - shutdown
			}
		}(worker)
	}

	log.Println("Server started!")

	err = http.ListenAndServe(conf.HTTP.Addr, handler)
	if err != nil {
		panic(err)
	}
}
