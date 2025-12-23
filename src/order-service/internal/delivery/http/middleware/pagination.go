package middleware

import (
	"context"
	"net/http"
	"strconv"
)

func PaginationMiddleware(page, limit int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if p := r.URL.Query().Get("page"); p != "" {
				if i, err := strconv.Atoi(p); err == nil && i > 0 {
					page = i
				}
			}

			if l := r.URL.Query().Get("limit"); l != "" {
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

func PageFromContext(ctx context.Context) int {
	page, ok := ctx.Value(ctxKeyPage).(int)
	if !ok {
		panic("page missing from context")
	}

	return page
}

func LimitFromContext(ctx context.Context) int {
	limit, ok := ctx.Value(ctxKeyLimit).(int)
	if !ok {
		panic("limit missing from context")
	}

	return limit
}
