package middleware

import (
	"context"
	"net/http"
	"orders-service/internal/delivery/http/helpers"

	"github.com/google/uuid"
)

func UserIdAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userIdStr := r.Header.Get("X-User-ID")

		if userIdStr == "" {
			helpers.RespondError(w, http.StatusUnauthorized, "missing X-User-ID")
			return
		}

		id, err := uuid.Parse(userIdStr)
		if err != nil {
			helpers.RespondError(w, http.StatusBadRequest, "Invalid X-User-ID UUID: "+userIdStr)
			return
		}

		ctx := context.WithValue(r.Context(), ctxKeyUserId, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserIdFromContext(ctx context.Context) uuid.UUID {
	id, ok := ctx.Value(ctxKeyUserId).(uuid.UUID)
	if !ok {
		panic("user id missing from context")
	}

	return id
}
