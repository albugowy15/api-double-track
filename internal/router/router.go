package router

import (
	"github.com/albugowy15/api-double-track/internal/controllers"
	userMiddleware "github.com/albugowy15/api-double-track/internal/middleware"
	"github.com/albugowy15/api-double-track/pkg/auth"
	"github.com/albugowy15/api-double-track/pkg/swagger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
)

func Init() *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Use(middleware.RequestID)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Compress(5, "text/html", "application/json"))

	router.Get("/swagger/*", swagger.WrapHandler)

	router.Route("/v1", func(r chi.Router) {
		r.Post("/auth/login", controllers.HandlePostLogin)
		r.Get("/alternatives", controllers.HandleGetAlternatives)

		// protected route Group
		// this route group require authentication and authorization
		r.Group(func(r chi.Router) {
			r.Use(jwtauth.Verifier(auth.GetAuth()))
			r.Use(auth.Authenticator)

			r.Get("/school", controllers.HandleGetSchool)

			// admin route
			r.Group(func(r chi.Router) {
				r.Use(userMiddleware.CheckAdminRole)
				r.Get("/admin/profile", controllers.HandleGetAdminProfile)
				r.Patch("/admin/profile", controllers.HandlePatchAdminProfile)
				r.Patch("/admin/change-password", controllers.HandlePatchAdminChangePassword)
				r.Get("/statistics", controllers.HandleGetStatistics)
				r.Post("/questionnare/settings", controllers.HandlePostQuestionnareSettings)
				r.Get("/questionnare/settings", controllers.HandleGetQuestionnareSettings)
				r.Get("/questionnare/settings/incomplete", controllers.HandleGetIncompleteQuestionnareSettings)
				r.Delete("/recommendations", controllers.HandleDeleteRecommendations)
				r.Get("/recommendations", controllers.HandleGetRecommendations)
				r.Get("/recommendations/student/{studentId}", controllers.HandleGetRecommendationStudent)
				r.Get("/students", controllers.HandleGetStudents)
				r.Post("/students", controllers.HandlePostStudent)
				r.Get("/students/{studentId}", controllers.HandleGetStudent)
				r.Patch("/students/{studentId}", controllers.HandlePatchStudent)
				r.Delete("/students", controllers.HandleDeleteStudent)
			})

			// student route
			r.Group(func(r chi.Router) {
				r.Use(userMiddleware.CheckStudentRole)
				r.Get("/students/profile", controllers.HandleGetStudentProfile)
				r.Patch("/students/profile", controllers.HandlePatchStudentProfile)
				r.Patch("/students/change-password", controllers.HandlePatchStudentChangePassword)
				r.Get("/recommendations/student", controllers.HandleGetRecommendationsStudent)
				r.Get("/questionnare/questions", controllers.HandleGetQuestions)
				r.Get("/questionnare/status", controllers.HandleGetQuesionnareStatus)
				r.Post("/questionnare/answers", controllers.HandlePostAnswers)
				r.Delete("/questionnare/answers", controllers.HandleDeleteAnswer)
				r.Get("/expectations", controllers.HandleGetExpectations)
				r.Get("/expectations/status", controllers.HandleGetExpectationsStatus)
				r.Post("/expectations", controllers.HandlePostExpectations)
				r.Patch("/expectations", controllers.HandlePatchExpectations)
				r.Delete("/expectations", controllers.HandleDeleteExpectations)
			})
		})
	})

	return router
}
