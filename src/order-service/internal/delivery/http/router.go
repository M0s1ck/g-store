package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"orders-service/internal/delivery/http/handlers"
	mymiddleware "orders-service/internal/delivery/http/middleware"
)

type RouterDeps struct {
	OrderHandler *handlers.OrderHandler
}

func NewRouter(deps *RouterDeps) http.Handler {
	router := chi.NewRouter()

	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
	)

	router.Route("/api/orders/{id}", func(r chi.Router) {
		r.Use(mymiddleware.UUIDMiddleware("id"))
		r.Get("/", deps.OrderHandler.GetById)
	})

	addHello(router)

	return router
}

func addHello(r *chi.Mux) {
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("Hello World"))
		if err != nil {
			return
		}
	})
}
