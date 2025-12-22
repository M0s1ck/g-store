package middleware

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func UUIDMiddleware(param string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw := chi.URLParam(r, param)

			id, err := uuid.Parse(raw)
			if err != nil {
				http.Error(w, "invalid uuid", http.StatusBadRequest)
				return
			}

			ctx := context.WithValue(r.Context(), ctxKeyUUID, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
