package controllers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/albugowy15/api-double-track/internal/pkg/models"
	"github.com/albugowy15/api-double-track/internal/pkg/repositories"
	"github.com/albugowy15/api-double-track/internal/pkg/schemas"
	"github.com/albugowy15/api-double-track/internal/pkg/utils/httputil"
	"github.com/albugowy15/api-double-track/internal/pkg/utils/jwt"
	"github.com/guregu/null/v5"
)

// GetRecommendations godoc
//
//	@Summary		Get all students recommendations
//	@Description	Get all students recommendations
//	@Tags			Recommendations
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		200				{object}	httputil.DataJsonResponse{data=[]schemas.StudentRecommendation}
//	@Failure		400				{object}	httputil.ErrorJsonResponse
//	@Failure		500				{object}	httputil.ErrorJsonResponse
//	@Router			/recommendations [get]
func GetRecommendations(w http.ResponseWriter, r *http.Request) {
	schoolIdClaim, _ := jwt.GetJwtClaim(r, "school_id")
	schoolId := schoolIdClaim.(string)
	studentRecommendations, err := repositories.GetRecommendationRepository().GetRecommendationsBySchoolId(schoolId)
	if err != nil {
		httputil.SendError(w, httputil.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	httputil.SendData(w, studentRecommendations, http.StatusOK)
}

// GetRecommendations godoc
//
//	@Summary		Get recommendations for a student
//	@Description	Get recommendations for a student
//	@Tags			Recommendations
//	@Tags			Student
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		200				{object}	httputil.DataJsonResponse{data=schemas.Recommendation}
//	@Failure		400				{object}	httputil.ErrorJsonResponse
//	@Failure		404				{object}	httputil.ErrorJsonResponse
//	@Failure		500				{object}	httputil.ErrorJsonResponse
//	@Router			/recommendations/student [get]
func GetStudentRecommendations(w http.ResponseWriter, r *http.Request) {
	studentIdClaim, _ := jwt.GetJwtClaim(r, "user_id")
	studentId := studentIdClaim.(string)

	consistencyRatio, err := repositories.GetRecommendationRepository().GetAHPConsistencyRatio(studentId)
	if err != nil {
		if err == sql.ErrNoRows {
			httputil.SendError(w, errors.New("isian kuesioner tidak ditemukan"), http.StatusNotFound)
			return
		}
		httputil.SendError(w, httputil.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	ahpResults, err := repositories.GetRecommendationRepository().GetAHPRecommendations(studentId)
	if err != nil {
		httputil.SendError(w, httputil.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	ahp := models.AhpRecommendation{
		Result:           ahpResults,
		ConsistencyRatio: null.FloatFrom(float64(consistencyRatio)),
	}

	res := models.Recommendation{
		Ahp: ahp,
	}
	httputil.SendData(w, res, http.StatusOK)
}

// DeleteRecommendations godoc
//
//	@Summary		Delete recommendation for a student
//	@Description	Delete recommendation for a student
//	@Tags			Recommendations
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string								true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			body			body		schemas.DeleteRecommendationRequest	true	"Delete student recommendation request body"
//	@Success		201				{object}	httputil.MessageJsonResponse
//	@Failure		400				{object}	httputil.ErrorJsonResponse
//	@Failure		500				{object}	httputil.ErrorJsonResponse
//	@Router			/recommendations [delete]
func DeleteRecommendations(w http.ResponseWriter, r *http.Request) {
	var body schemas.DeleteRecommendationRequest
	httputil.GetBody(w, r, &body)
	if len(body.StudentId) == 0 {
		httputil.SendError(w, errors.New("id siswa wajib diisi"), http.StatusBadRequest)
		return
	}

	err := repositories.GetAnswersRepository().DeleteAnswers(body.StudentId)
	if err != nil {
		log.Println(err)
		httputil.SendError(w, httputil.ErrInternalServer, http.StatusBadRequest)
		return
	}

	httputil.SendMessage(w, "berhasil menghapus rekomendasi", http.StatusCreated)
}

// GetStudentRecommendationDetail godoc
//
//	@Summary		Get student recommendations details
//	@Description	Get student recommendations details
//	@Tags			Recommendations
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			studentId		path		string	true	"Id student"
//	@Success		200				{object}	httputil.DataJsonResponse{data=schemas.Recommendation}
//	@Failure		400				{object}	httputil.ErrorJsonResponse
//	@Failure		404				{object}	httputil.ErrorJsonResponse
//	@Failure		500				{object}	httputil.ErrorJsonResponse
//	@Router			/recommendations/student/{studentId} [get]
func GetStudentRecommendationDetail(w http.ResponseWriter, r *http.Request) {
	studentId := r.PathValue("studentId")
	if len(studentId) == 0 {
		httputil.SendError(w, errors.New("id siswa wajib diisi"), http.StatusBadRequest)
		return
	}
	consistencyRatio, err := repositories.GetRecommendationRepository().GetAHPConsistencyRatio(studentId)
	if err != nil {
		if err == sql.ErrNoRows {
			httputil.SendError(w, errors.New("isian kuesioner tidak ditemukan"), http.StatusNotFound)
			return
		}
		httputil.SendError(w, httputil.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	ahpResults, err := repositories.GetRecommendationRepository().GetAHPRecommendations(studentId)
	if err != nil {
		httputil.SendError(w, httputil.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	ahp := models.AhpRecommendation{
		Result:           ahpResults,
		ConsistencyRatio: null.FloatFrom(float64(consistencyRatio)),
	}

	res := models.Recommendation{
		Ahp: ahp,
	}
	httputil.SendData(w, res, http.StatusOK)
}
