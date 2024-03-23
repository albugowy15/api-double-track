package controllers

import (
	"log"
	"net/http"

	"github.com/albugowy15/api-double-track/internal/pkg/models"
	"github.com/albugowy15/api-double-track/internal/pkg/repositories"
	"github.com/albugowy15/api-double-track/internal/pkg/utils/httputil"
	"github.com/albugowy15/api-double-track/internal/pkg/utils/jwt"
)

// GetStatistics godoc
//
//	@Summary		Get statistic
//	@Description	Get statistics
//	@Tags			Common
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		200				{object}	httputil.DataJsonResponse{data=schemas.Statistic}
//	@Failure		500				{object}	httputil.ErrorJsonResponse
//	@Router			/statistics [get]
func GetStatistics(w http.ResponseWriter, r *http.Request) {
	res := models.Statistic{
		RegisteredStudents:       123,
		QuestionnareCompleted:    450,
		RecommendationAcceptance: 90.34,
		ConsistencyAvg:           92.54,
	}
	httputil.SendData(w, res, http.StatusOK)
}

// GetAlternatives godoc
//
//	@Summary		Get alternatives
//	@Description	Get all alternatives
//	@Tags			Common
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	httputil.DataJsonResponse{data=[]schemas.Alternative}
//	@Failure		500	{object}	httputil.ErrorJsonResponse
//	@Router			/alternatives [get]
func GetAlternatives(w http.ResponseWriter, r *http.Request) {
	s := repositories.GetAlternativeRepository()
	alternatives, err := s.GetAlternatives()
	if err != nil {
		httputil.SendError(w, httputil.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	httputil.SendData(w, alternatives, http.StatusOK)
}

// GetSchool godoc
//
//	@Summary		Get school
//	@Description	Get current authenticated user shcool
//	@Tags			Common
//	@Tags			Student
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		200				{object}	httputil.DataJsonResponse{data=schemas.School}
//	@Failure		500				{object}	httputil.ErrorJsonResponse
//	@Router			/school [get]
func GetSchool(w http.ResponseWriter, r *http.Request) {
	// get school_id from token
	schoolIdClaim, _ := jwt.GetJwtClaim(r, "school_id")
	schoolId := schoolIdClaim.(string)

	school, err := repositories.GetSchoolRepository().GetSchoolById(schoolId)
	if err != nil {
		log.Printf("err get school: %v", err)
		httputil.SendError(w, httputil.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	httputil.SendData(w, school, http.StatusOK)
}
