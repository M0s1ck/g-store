package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/swaggo/http-swagger"

	"orders-service/internal/delivery/http/handlers"
	mymiddleware "orders-service/internal/delivery/http/middleware"
	"orders-service/internal/docs"
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

	router.Route("/v1/orders", func(r chi.Router) {
		r.Use(mymiddleware.UserIdAuthMiddleware)

		r.With(mymiddleware.PaginationMiddleware(1, 20)).
			Get("/", deps.OrderHandler.GetByUser)

		r.Post("/", deps.OrderHandler.Create)

		r.Route("/{id}", func(r chi.Router) {
			r.Use(mymiddleware.UUIDMiddleware("id"))
			r.Get("/", deps.OrderHandler.GetById)
		})
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

func addSwagger(r chi.Router) {
	r.Get("/swagger/*", func(w http.ResponseWriter, req *http.Request) {
		// Подставляем host динамически
		docs.SwaggerInfo.Host = req.Host
		docs.SwaggerInfo.Schemes = []string{"http"}

		// Учитываем gateway prefix
		if p := req.Header.Get("X-Forwarded-Prefix"); p != "" {
			docs.SwaggerInfo.BasePath = p + "/v1"
		} else {
			docs.SwaggerInfo.BasePath = "/v1"
		}

		httpSwagger.WrapHandler(w, req)
	})
}
