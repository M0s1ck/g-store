package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	"payment-service/internal/delivery/http/dto"
	"payment-service/internal/delivery/http/handlers"
	mymiddleware "payment-service/internal/delivery/http/middleware"
	"payment-service/internal/docs"
)

type RouterDeps struct {
	AccountHandler *handlers.AccountHandler
}

func NewRouter(deps *RouterDeps) http.Handler {
	router := chi.NewRouter()

	router.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
	)

	router.Route("/v1/accounts", func(r chi.Router) {
		r.Use(mymiddleware.UserIdAuthMiddleware)

		r.With(mymiddleware.UUIDMiddleware("id")).
			Get("/{id}", deps.AccountHandler.GetById)

		r.Post("/", deps.AccountHandler.Create)

		r.With(mymiddleware.UUIDMiddleware("id")).
			With(mymiddleware.BindJSONBodyMiddleware[dto.TopUpReq]()).
			Patch("/{id}", deps.AccountHandler.TopUp)
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
