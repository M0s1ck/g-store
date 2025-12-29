package my_middleware

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

func BindJSONBodyMiddleware[T any]() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var body T

			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				http.Error(w, "invalid JSON body", http.StatusBadRequest)
				return
			}

			defer func() {
				_ = r.Body.Close()
			}()

			ctx := context.WithValue(r.Context(), ctxKeyBody, body)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func BodyFromContext[T any](ctx context.Context) T {
	body, ok := ctx.Value(ctxKeyBody).(T)
	if !ok {
		panic("missing body from context")
	}

	return body
}

func BindOptionalJSONBodyMiddleware[T any]() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			var body *T

			if r.Body != nil {
				defer func() { _ = r.Body.Close() }()

				decoder := json.NewDecoder(r.Body)
				decoder.DisallowUnknownFields()

				var tmp T
				if err := decoder.Decode(&tmp); err != nil {
					// EOF = пустое тело → это ок
					if !errors.Is(err, io.EOF) {
						http.Error(w, "invalid JSON body", http.StatusBadRequest)
						return
					}
				} else {
					body = &tmp
				}
			}

			ctx := context.WithValue(r.Context(), ctxKeyBody, body)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func OptionalBodyFromContext[T any](ctx context.Context) *T {
	body := ctx.Value(ctxKeyBody).(*T)
	return body
}
