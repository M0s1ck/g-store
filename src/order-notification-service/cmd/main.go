package main

import (
	"log"

	"order-notification-service/internal/app"
	"order-notification-service/internal/config"
)

func main() {
	log.Println("Service is starting...")

	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	wsApp := app.Build(conf)
	wsApp.Run()
}
