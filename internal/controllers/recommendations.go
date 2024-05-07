package controllers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/albugowy15/api-double-track/internal/models"
	"github.com/albugowy15/api-double-track/internal/repositories"
	"github.com/albugowy15/api-double-track/internal/utils"
	"github.com/albugowy15/api-double-track/pkg/auth"
	"github.com/albugowy15/api-double-track/pkg/httpx"
	"github.com/albugowy15/api-double-track/pkg/schemas"
	"github.com/guregu/null/v5"
)

// HandleGetRecommendations godoc
//
//	@Summary		Get all students recommendations
//	@Description	Get all students recommendations
//	@Tags			Recommendations
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		200				{object}	httpx.DataJsonResponse{data=[]schemas.StudentRecommendation}
//	@Failure		400				{object}	httpx.ErrorJsonResponse
//	@Failure		500				{object}	httpx.ErrorJsonResponse
//	@Router			/recommendations [get]
func HandleGetRecommendations(w http.ResponseWriter, r *http.Request) {
	schoolIdClaim, _ := auth.GetJwtClaim(r, "school_id")
	schoolId := schoolIdClaim.(string)
	studentRecommendations, err := repositories.GetRecommendationsBySchoolId(schoolId)
	if err != nil {
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	httpx.SendData(w, studentRecommendations, http.StatusOK)
}

// HandleGetRecommendationsStudent godoc
//
//	@Summary		Get recommendations for a student
//	@Description	Get recommendations for a student
//	@Tags			Recommendations
//	@Tags			Student
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		200				{object}	httpx.DataJsonResponse{data=schemas.Recommendation}
//	@Failure		400				{object}	httpx.ErrorJsonResponse
//	@Failure		404				{object}	httpx.ErrorJsonResponse
//	@Failure		500				{object}	httpx.ErrorJsonResponse
//	@Router			/recommendations/student [get]
func HandleGetRecommendationsStudent(w http.ResponseWriter, r *http.Request) {
	studentIdClaim, _ := auth.GetJwtClaim(r, "user_id")
	studentId := studentIdClaim.(string)

	consistencyRatio, err := repositories.GetAHPConsistencyRatio(studentId)
	if err != nil {
		if err == sql.ErrNoRows {
			httpx.SendError(w, errors.New("isian kuesioner tidak ditemukan"), http.StatusNotFound)
			return
		}
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	ahpResults, err := repositories.GetAHPRecommendations(studentId)
	if err != nil {
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	ahpResultsWithRanks, err := utils.MakeRecommendationRanks(ahpResults)
	if err != nil {
		httpx.SendError(w, err, http.StatusBadRequest)
		return
	}
	ahp := models.AhpRecommendation{
		Result:           ahpResultsWithRanks,
		ConsistencyRatio: null.FloatFrom(float64(consistencyRatio)),
	}

	topsisResults, err := repositories.GetTOPSISRecommendations(studentId)
	if err != nil {
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	topsis := models.TopsisRecommendation{
		Result: topsisResults,
	}

	topsisAHPResults, err := repositories.GetTOPSISAHPRecommendations(studentId)
	if err != nil {
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	topsis_ahp := models.TopsisAHPRecommendation{
		Result: topsisAHPResults,
	}
	res := models.Recommendation{
		Ahp:       ahp,
		Topsis:    topsis,
		TopsisAHP: topsis_ahp,
	}

	httpx.SendData(w, res, http.StatusOK)
}

// HandleDeleteRecommendations godoc
//
//	@Summary		Delete recommendation for a student
//	@Description	Delete recommendation for a student
//	@Tags			Recommendations
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string								true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			body			body		schemas.DeleteRecommendationRequest	true	"Delete student recommendation request body"
//	@Success		201				{object}	httpx.MessageJsonResponse
//	@Failure		400				{object}	httpx.ErrorJsonResponse
//	@Failure		500				{object}	httpx.ErrorJsonResponse
//	@Router			/recommendations [delete]
func HandleDeleteRecommendations(w http.ResponseWriter, r *http.Request) {
	var body schemas.DeleteRecommendationRequest
	if err := httpx.GetBody(r, &body); err != nil {
		httpx.SendError(w, err, http.StatusBadRequest)
	}
	if len(body.StudentId) == 0 {
		httpx.SendError(w, errors.New("id siswa wajib diisi"), http.StatusBadRequest)
		return
	}

	err := repositories.DeleteAnswers(body.StudentId)
	if err != nil {
		log.Println(err)
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusBadRequest)
		return
	}

	httpx.SendMessage(w, "berhasil menghapus rekomendasi", http.StatusCreated)
}

// HandleGetRecommendationStudent godoc
//
//	@Summary		Get student recommendations details
//	@Description	Get student recommendations details
//	@Tags			Recommendations
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			studentId		path		string	true	"Id student"
//	@Success		200				{object}	httpx.DataJsonResponse{data=schemas.Recommendation}
//	@Failure		400				{object}	httpx.ErrorJsonResponse
//	@Failure		404				{object}	httpx.ErrorJsonResponse
//	@Failure		500				{object}	httpx.ErrorJsonResponse
//	@Router			/recommendations/student/{studentId} [get]
func HandleGetRecommendationStudent(w http.ResponseWriter, r *http.Request) {
	studentId := r.PathValue("studentId")
	if len(studentId) == 0 {
		httpx.SendError(w, errors.New("id siswa wajib diisi"), http.StatusBadRequest)
		return
	}
	consistencyRatio, err := repositories.GetAHPConsistencyRatio(studentId)
	if err != nil {
		if err == sql.ErrNoRows {
			httpx.SendError(w, errors.New("isian kuesioner tidak ditemukan"), http.StatusNotFound)
			return
		}
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	ahpResults, err := repositories.GetAHPRecommendations(studentId)
	if err != nil {
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	ahpResultsWithRanks, err := utils.MakeRecommendationRanks(ahpResults)
	if err != nil {
		httpx.SendError(w, err, http.StatusBadRequest)
		return
	}
	ahp := models.AhpRecommendation{
		Result:           ahpResultsWithRanks,
		ConsistencyRatio: null.FloatFrom(float64(consistencyRatio)),
	}

	topsisResults, err := repositories.GetTOPSISRecommendations(studentId)
	if err != nil {
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	topsis := models.TopsisRecommendation{
		Result: topsisResults,
	}

	topsisAHPResults, err := repositories.GetTOPSISAHPRecommendations(studentId)
	if err != nil {
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	topsis_ahp := models.TopsisAHPRecommendation{
		Result: topsisAHPResults,
	}

	res := models.Recommendation{
		Ahp:       ahp,
		Topsis:    topsis,
		TopsisAHP: topsis_ahp,
	}
	// fmt.Println("result : ", res.Topsis)
	httpx.SendData(w, res, http.StatusOK)
}
