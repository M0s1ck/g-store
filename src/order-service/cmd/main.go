package main

import (
	"log"
	"net/http"

	"orders-service/internal/app"
	"orders-service/internal/config"
)

func main() {
	conf, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	handler := app.Build(conf)

	log.Println("Server started!")

	err = http.ListenAndServe(":8080", *handler)
	if err != nil {
		panic(err)
	}
}
