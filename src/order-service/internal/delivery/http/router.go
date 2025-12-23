package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/swaggo/http-swagger"

	"orders-service/internal/delivery/http/handlers"
	mymiddleware "orders-service/internal/delivery/http/middleware"
	_ "orders-service/internal/docs"
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

	router.Route("/v1/orders/{id}", func(r chi.Router) {
		r.Use(mymiddleware.UUIDMiddleware("id"))
		r.Use(mymiddleware.UserIdAuthMiddleware)
		r.Get("/", deps.OrderHandler.GetById)
	})

	addHello(router)
	addSwagger(router)

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

func addSwagger(r *chi.Mux) {
	r.Get("/swagger/*", httpSwagger.WrapHandler)
}
