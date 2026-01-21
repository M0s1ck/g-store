package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"github.com/M0s1ck/g-store/src/pkg/http/responds"
)

func UUIDMiddleware(param string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw := chi.URLParam(r, param)

			id, err := uuid.Parse(raw)
			if err != nil {
				responds.RespondError(w, http.StatusBadRequest, errors.New("Invalid UUID: "+raw))
				return
			}

			ctx := context.WithValue(r.Context(), ctxKeyUUID, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UUIDFromContext(ctx context.Context) (uuid.UUID, error) {
	id, ok := ctx.Value(ctxKeyUUID).(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("uuid missing from context")
	}

	return id, nil
}
