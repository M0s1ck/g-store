package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"

	mymiddleware "orders-service/internal/delivery/http/middleware"
)

type OrderHandler struct {
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{}
}

func (h *OrderHandler) GetById(w http.ResponseWriter, r *http.Request) {
	orderId := mymiddleware.UUIDFromContext(r.Context())

	hello := Hello{
		Id:   orderId,
		Name: "Hello World",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(hello)
	if err != nil {
		panic(err)
	}
}

type Hello struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
