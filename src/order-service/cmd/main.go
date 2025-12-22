package main

import (
	"net/http"

	"orders-service/internal/app"
	"orders-service/internal/config"
)

func main() {
	conf := config.NewConfig()

	handler := app.Build(conf)

	err := http.ListenAndServe(":8080", *handler)
	if err != nil {
		panic(err)
	}
}
