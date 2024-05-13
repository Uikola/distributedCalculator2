package http

import (
	"github.com/riandyrn/otelchi"
	"net/http"

	"github.com/Uikola/distributedCalculator2/orchestrator/internal/server/http/cresource"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/server/http/expression"
	"github.com/Uikola/distributedCalculator2/orchestrator/internal/server/http/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func NewServer(
	userHandler *user.Handler,
	expressionHandler *expression.Handler,
	cResourceHandler *cresource.Handler,
) http.Handler {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	router.Use(otelchi.Middleware("orchestrator", otelchi.WithChiRoutes(router)))
	addRoutes(router, userHandler, expressionHandler, cResourceHandler)

	var handler http.Handler = router

	return handler
}

func addRoutes(
	router *chi.Mux,
	userHandler *user.Handler,
	expressionHandler *expression.Handler,
	cResourceHandler *cresource.Handler,
) {

	router.Route("/api", func(r chi.Router) {
		r.Post("/register", userHandler.Register)
		r.Post("/login", userHandler.Login)

		r.With(userHandler.Auth).Group(func(r chi.Router) {
			r.Post("/calculate", expressionHandler.AddExpression)
			r.Get("/expressions/{id}", expressionHandler.GetExpression)
			r.Get("/expressions", expressionHandler.ListExpressions)
			r.Get("/results/{id}", expressionHandler.GetResult)
			r.Get("/cresources", cResourceHandler.ListCResources)
			r.Get("/operations", userHandler.ListOperations)
			r.Put("/operations", userHandler.UpdateOperation)
		})
	})
}
