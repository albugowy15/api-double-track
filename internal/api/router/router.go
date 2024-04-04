package router

import (
	"github.com/albugowy15/api-double-track/internal/api/controllers"
	userMiddleware "github.com/albugowy15/api-double-track/internal/api/middleware"
	"github.com/albugowy15/api-double-track/internal/pkg/swagger"
	"github.com/albugowy15/api-double-track/internal/pkg/utils/jwt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
)

func Setup() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Compress(5, "text/html", "application/json"))

	router.Get("/swagger/*", swagger.WrapHandler)

	router.Route("/v1", func(r chi.Router) {
		r.Post("/auth/login", controllers.Login)
		r.Get("/alternatives", controllers.GetAlternatives)

		// protected route Group
		// this route group require authentication and authorization
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(jwt.GetAuth()))
			r.Use(jwt.Authenticator)

			r.Get("/school", controllers.GetSchool)

			// admin route
			r.Group(func(r chi.Router) {
				r.Use(userMiddleware.CheckAdminRole)
				r.Get("/admin/profile", controllers.GetAdminProfile)
				r.Patch("/admin/profile", controllers.UpdateAdminProfile)
				r.Get("/statistics", controllers.GetStatistics)
				r.Post("/questionnare/settings", controllers.AddQuestionnareSettings)
				r.Get("/questionnare/settings", controllers.GetQuestionnareSettings)
				r.Get("/questionnare/settings/incomplete", controllers.GetIncompleteQuestionnareSettings)
				r.Delete("/recommendations", controllers.DeleteRecommendations)
				r.Get("/recommendations", controllers.GetRecommendations)
				r.Get("/recommendations/student/{studentId}", controllers.GetStudentRecommendationDetail)
				r.Get("/students", controllers.GetStudents)
				r.Post("/students", controllers.AddStudent)
				r.Get("/students/{studentId}", controllers.GetStudent)
				r.Patch("/students/{studentId}", controllers.UpdateStudent)
				r.Delete("/students", controllers.DeleteStudent)
			})

			// student route
			r.Group(func(r chi.Router) {
				r.Use(userMiddleware.CheckStudentRole)
				r.Get("/students/profile", controllers.GetStudentProfile)
				r.Patch("/students/profile", controllers.UpdateStudentProfile)
				r.Get("/recommendations/student", controllers.GetStudentRecommendations)
				r.Get("/questionnare/questions", controllers.GetQuestions)
				r.Get("/questionnare/status", controllers.GetQuesionnareStatus)
				r.Post("/questionnare/answers", controllers.SubmitAnswer)
				r.Delete("/questionnare/answers", controllers.DeleteAnswer)
			})
		})
	})

	return router
}
