package controllers

import (
	"errors"
	"net/http"

	"github.com/albugowy15/api-double-track/internal/models"
	"github.com/albugowy15/api-double-track/internal/repositories"
	"github.com/albugowy15/api-double-track/internal/validator"
	"github.com/albugowy15/api-double-track/pkg/auth"
	"github.com/albugowy15/api-double-track/pkg/httpx"
)

func HandleGetExpectations(w http.ResponseWriter, r *http.Request) {
	studentIdClaim, _ := auth.GetJwtClaim(r, "user_id")
	studentId := studentIdClaim.(string)

	expectations, err := repositories.GetStudentExpectations(studentId)
	if err != nil {
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	httpx.SendData(w, expectations, http.StatusOK)
}

func HandleGetExpectationsStatus(w http.ResponseWriter, r *http.Request) {
	studentIdClaim, _ := auth.GetJwtClaim(r, "user_id")
	studentId := studentIdClaim.(string)

	expectationExist := repositories.CheckStudentExpectationExist(studentId)
	if expectationExist {
		res := models.ExpectationStatusResponse{
			Status: "COMPLETED",
		}
		httpx.SendData(w, res, http.StatusOK)
		return
	}
	res := models.ExpectationStatusResponse{
		Status: "READY",
	}
	httpx.SendData(w, res, http.StatusOK)
}

func HandlePostExpectations(w http.ResponseWriter, r *http.Request) {
	studentIdClaim, _ := auth.GetJwtClaim(r, "user_id")
	studentId := studentIdClaim.(string)
	var body models.ExpectationRequest
	if err := httpx.GetBody(r, &body); err != nil {
		httpx.SendError(w, httpx.ErrDecodeJsonBody, http.StatusBadRequest)
		return
	}

	if err := validator.ValidateExpectationRequest(body); err != nil {
		httpx.SendError(w, err, http.StatusBadRequest)
		return
	}
	if repositories.CheckStudentExpectationExist(studentId) {
		httpx.SendError(w, errors.New("siswa telah mengisi ekspektasi rekomendasi"), http.StatusBadRequest)
		return
	}
	//
	//
	// save to db
	err := repositories.SaveStudentExpectations(body.Expectations, studentId)
	if err != nil {
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	httpx.SendMessage(w, "berhasil menyimpan ekspektasi rekomendasi", http.StatusCreated)
}

func HandlePatchExpectations(w http.ResponseWriter, r *http.Request) {
	studentIdClaim, _ := auth.GetJwtClaim(r, "user_id")
	studentId := studentIdClaim.(string)
	var body models.ExpectationRequest
	if err := httpx.GetBody(r, &body); err != nil {
		httpx.SendError(w, httpx.ErrDecodeJsonBody, http.StatusBadRequest)
		return
	}

	if err := validator.ValidateExpectationRequest(body); err != nil {
		httpx.SendError(w, err, http.StatusBadRequest)
		return
	}
	if !repositories.CheckStudentExpectationExist(studentId) {
		httpx.SendError(w, errors.New("siswa belum mengisi ekspektasi rekomendasi"), http.StatusBadRequest)
		return
	}
	//
	//
	// save to db
	err := repositories.UpdateStudentExpectation(body.Expectations, studentId)
	if err != nil {
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	httpx.SendMessage(w, "berhasil memperbarui ekspektasi rekomendasi", http.StatusCreated)
}

func HandleDeleteExpectations(w http.ResponseWriter, r *http.Request) {
	studentIdClaim, _ := auth.GetJwtClaim(r, "user_id")
	studentId := studentIdClaim.(string)

	if err := repositories.DeleteStudentExpectation(studentId); err != nil {
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	httpx.SendMessage(w, "berhasil menghapus ekspektasi rekomendasi", http.StatusOK)
}
