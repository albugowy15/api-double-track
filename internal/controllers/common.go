package controllers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/albugowy15/api-double-track/internal/models"
	"github.com/albugowy15/api-double-track/internal/repositories"
	"github.com/albugowy15/api-double-track/pkg/auth"
	"github.com/albugowy15/api-double-track/pkg/httpx"
)

// HandleGetStatistics godoc
//
//	@Summary		Get statistic
//	@Description	Get statistics
//	@Tags			Common
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		200				{object}	httpx.DataJsonResponse{data=schemas.Statistic}
//	@Failure		500				{object}	httpx.ErrorJsonResponse
//	@Router			/statistics [get]
func HandleGetStatistics(w http.ResponseWriter, r *http.Request) {
	res := models.Statistic{
		RegisteredStudents:       123,
		QuestionnareCompleted:    450,
		RecommendationAcceptance: 90.34,
		ConsistencyAvg:           92.54,
	}
	httpx.SendData(w, res, http.StatusOK)
}

// HandleGetAlternatives godoc
//
//	@Summary		Get alternatives
//	@Description	Get all alternatives
//	@Tags			Common
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	httpx.DataJsonResponse{data=[]schemas.Alternative}
//	@Failure		500	{object}	httpx.ErrorJsonResponse
//	@Router			/alternatives [get]
func HandleGetAlternatives(w http.ResponseWriter, r *http.Request) {
	s := repositories.GetAlternativeRepository()
	alternatives, err := s.GetAlternatives()
	if err != nil {
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	httpx.SendData(w, alternatives, http.StatusOK)
}

// HandleGetSchool godoc
//
//	@Summary		Get school
//	@Description	Get current authenticated user shcool
//	@Tags			Common
//	@Tags			Student
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		200				{object}	httpx.DataJsonResponse{data=schemas.School}
//	@Failure		404				{object}	httpx.ErrorJsonResponse
//	@Failure		500				{object}	httpx.ErrorJsonResponse
//	@Router			/school [get]
func HandleGetSchool(w http.ResponseWriter, r *http.Request) {
	// get school_id from token
	schoolIdClaim, _ := auth.GetJwtClaim(r, "school_id")
	schoolId := schoolIdClaim.(string)

	school, err := repositories.GetSchoolRepository().GetSchoolById(schoolId)
	if err != nil {
		if err == sql.ErrNoRows {
			httpx.SendError(w, errors.New("sekolah tidak ditemukan"), http.StatusNotFound)
			return
		}
		log.Printf("err get school: %v", err)
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	httpx.SendData(w, school, http.StatusOK)
}
