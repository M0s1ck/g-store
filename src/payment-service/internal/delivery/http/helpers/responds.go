package helpers

import (
	"encoding/json"
	"log"
	"net/http"

	"payment-service/internal/delivery/http/dto"
)

func RespondJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		log.Printf("json encode error: %v", err)
	}
}

func RespondError(w http.ResponseWriter, status int, message string) {
	RespondJSON(w, status, dto.ErrorResponse{
		Error: message,
	})
}
