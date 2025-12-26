package main

import (
	"context"
	"log"
	"net/http"

	"payment-service/internal/app"
	"payment-service/internal/config"
	_ "payment-service/internal/docs"
	"payment-service/internal/infrastructure/background_workers"
)

// @title Payment Service API
// @version 1.0
// @description Payment microservice
// @BasePath /v1
// @schemes http https
func main() {
	log.Println("Service is starting...")

	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	handler, workers := app.Build(conf)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// launching background workers
	for _, worker := range workers {
		go func(worker background_workers.BackgroundWorker) {
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
		log.Printf("http server error: %v", err)
		cancel()
	}
}
