package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/M0s1ck/g-store/src/pkg/http/responds"
)

// BindJSONBodyMiddleware binds json body to the given type T.
// Responds with responds.ErrorResponse 400 if binding failed.
// Saves value to request context, use BodyFromContext to get it
func BindJSONBodyMiddleware[T any]() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				_ = r.Body.Close()
			}()

			var body T

			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				responds.RespondError(w, http.StatusBadRequest, errors.New("invalid JSON body"))
				return
			}

			ctx := context.WithValue(r.Context(), ctxKeyBody, body)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// BodyFromContext retrieves typed object from the context.
// It implies that BodyFromContext was used
func BodyFromContext[T any](ctx context.Context) (*T, error) {
	body, ok := ctx.Value(ctxKeyBody).(T)
	if !ok {
		return nil, errors.New("missing body from context")
	}

	return &body, nil
}
