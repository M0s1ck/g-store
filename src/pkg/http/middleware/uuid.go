package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"

	"github.com/M0s1ck/g-store/src/pkg/http/responds"
)

// UUIDMiddleware parses uuid from path key.
// Raw str is extracted by func of your framework (i.e. chi.URLParam(r, "id") or c.Param("id")).
// Responds with responds.ErrorResponse 400 if uuid parsing failed.
// Saves value to request context, use UUIDFromContext to get it
func UUIDMiddleware(extract func(*http.Request) string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			raw := extract(r)

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

// UUIDFromContext retrieves parsed uuid key from context.
// It implies that UUIDMiddleware was used
func UUIDFromContext(ctx context.Context) (uuid.UUID, error) {
	id, ok := ctx.Value(ctxKeyUUID).(uuid.UUID)
	if !ok {
		return uuid.Nil, errors.New("uuid missing from context")
	}

	return id, nil
}
