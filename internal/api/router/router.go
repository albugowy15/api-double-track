package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.con/albugowy15/api-double-track/internal/api/controllers"
)

func Setup() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Compress(5, "text/html", "application/json"))

	router.Route("/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Get("/hello", controllers.Hello)
		})
	})

	return router
}
