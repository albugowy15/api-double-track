package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"github.con/albugowy15/api-double-track/internal/api/controllers"
	userMiddleware "github.con/albugowy15/api-double-track/internal/api/middleware"
	"github.con/albugowy15/api-double-track/internal/pkg/utils/jwt"
)

func Setup() *chi.Mux {
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Compress(5, "text/html", "application/json"))

	router.Route("/v1", func(r chi.Router) {
		r.Post("/auth/login", controllers.Login)
		r.Post("/auth/register", func(w http.ResponseWriter, r *http.Request) {})
		r.Get("/alternatives", controllers.GetAlternatives)

		// protected route Group
		// this route group require authentication and authorization
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(jwt.GetAuth()))
			r.Use(jwt.Authenticator)

			// admin route
			r.Group(func(r chi.Router) {
				r.Use(userMiddleware.CheckAdminRole)
				r.Get("/statistics", controllers.GetStatistics)
				r.Post("/alternatives/settings", func(w http.ResponseWriter, r *http.Request) {})
				r.Delete("/recommendations", func(w http.ResponseWriter, r *http.Request) {})
				r.Get("/recommendations", func(w http.ResponseWriter, r *http.Request) {})
				r.Get("/students", controllers.GetStudents)
				r.Post("/students", controllers.AddStudent)
				r.Get("/students/{studentId}", controllers.GetStudent)
				r.Patch("/students/{studentId}", controllers.UpdateStudent)
				r.Delete("/students", controllers.DeleteStudent)
			})

			// student route
			r.Group(func(r chi.Router) {
				r.Use(userMiddleware.CheckStudentRole)
				r.Get("/recommendations/{studentId}", func(w http.ResponseWriter, r *http.Request) {})
				r.Get("/questionnare/questions", controllers.GetQuestions)
				r.Post("/questionnare/answers", controllers.SubmitAnswer)
			})
		})
	})

	return router
}
