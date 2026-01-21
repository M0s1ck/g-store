package middleware

import (
	"context"
	"errors"
	"net/http"
	"strconv"
)

// PaginationMiddleware parses query pagination params pageKey and limitKey.
// If params are not provided or invalid, uses default values page and limit.
// Saves values to request context, use PageFromContext and LimitFromContext to get them
func PaginationMiddleware(pageKey, limitKey string, page, limit int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if p := r.URL.Query().Get(pageKey); p != "" {
				if i, err := strconv.Atoi(p); err == nil && i > 0 {
					page = i
				}
			}

			if l := r.URL.Query().Get(limitKey); l != "" {
				if i, err := strconv.Atoi(l); err == nil && i > 0 {
					limit = i
				}
			}

			ctx := context.WithValue(r.Context(), ctxKeyPage, page)
			ctx = context.WithValue(ctx, ctxKeyLimit, limit)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// PageFromContext retrieves page value from context.
// It implies that PaginationMiddleware was used
func PageFromContext(ctx context.Context) (int, error) {
	page, ok := ctx.Value(ctxKeyPage).(int)
	if !ok {
		return 0, errors.New("page missing from context")
	}

	return page, nil
}

// LimitFromContext retrieves parsed limit value from context.
// It implies that PaginationMiddleware was used
func LimitFromContext(ctx context.Context) (int, error) {
	limit, ok := ctx.Value(ctxKeyLimit).(int)
	if !ok {
		return 0, errors.New("limit missing from context")
	}

	return limit, nil
}
