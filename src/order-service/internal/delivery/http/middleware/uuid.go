package middleware

import (
	"context"
	"net/http"
	"orders-service/internal/delivery/http/helpers"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func UUIDMiddleware(param string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw := chi.URLParam(r, param)

			id, err := uuid.Parse(raw)
			if err != nil {
				helpers.RespondError(w, http.StatusBadRequest, "Invalid UUID: "+raw)
				return
			}

			ctx := context.WithValue(r.Context(), ctxKeyUUID, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UUIDFromContext(ctx context.Context) uuid.UUID {
	id, ok := ctx.Value(ctxKeyUUID).(uuid.UUID)
	if !ok {
		panic("uuid missing from context")
	}

	return id
}
