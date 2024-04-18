package controllers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/albugowy15/api-double-track/internal/models"
	"github.com/albugowy15/api-double-track/internal/repositories"
	"github.com/albugowy15/api-double-track/internal/repositories/user"
	"github.com/albugowy15/api-double-track/internal/validator"
	"github.com/albugowy15/api-double-track/pkg/auth"
	"github.com/albugowy15/api-double-track/pkg/httpx"
	"github.com/albugowy15/api-double-track/pkg/schemas"
	"golang.org/x/crypto/bcrypt"
)

// HandlePostLogin godoc
//
//	@Summary		Login authentication
//	@Description	Login authentication for student and admin
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		schemas.LoginRequest	true	"Login request body"
//	@Success		200		{object}	httpx.DataJsonResponse{data=schemas.LoginResponse}
//	@Failure		400		{object}	httpx.ErrorJsonResponse
//	@Failure		500		{object}	httpx.ErrorJsonResponse
//	@Router			/auth/login [post]
func HandlePostLogin(w http.ResponseWriter, r *http.Request) {
	var body models.LoginRequest
	if err := httpx.GetBody(r, &body); err != nil {
		httpx.SendError(w, err, http.StatusBadRequest)
	}
	err := validator.ValidateLoginRequest(body)
	if err != nil {
		httpx.SendError(w, err, http.StatusBadRequest)
		return
	}

	switch body.Type {
	case "admin":
		admin, err := user.GetAdminByUsername(body.Username)
		if err != nil {
			httpx.SendError(w, errors.New("username atau password salah"), http.StatusBadRequest)
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(body.Password)); err != nil {
			httpx.SendError(w, errors.New("username atau password salah"), http.StatusBadRequest)
			return
		}
		school, err := repositories.GetSchoolByAdminId(admin.Id)
		if err != nil {
			httpx.SendError(w, errors.New("admin tidak memiliki akses ke sekolah"), http.StatusBadRequest)
			return
		}
		claim := auth.JWTClaim{
			UserId:   admin.Id,
			Username: admin.Username,
			Email:    admin.Email.String,
			Role:     "admin",
			SchoolId: school.Id,
		}
		token := auth.CreateToken(claim)
		res := models.LoginResponse{
			Token:    token,
			Username: admin.Username,
			Id:       admin.Id,
			Role:     "admin",
			SchoolId: school.Id,
		}
		httpx.SendData(w, res, http.StatusOK)
		return
	case "student":
		student, err := user.GetStudentByUsername(body.Username)
		if err != nil {
			log.Println(err)
			httpx.SendError(w, errors.New("username atau password salah"), http.StatusBadRequest)
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(body.Password)); err != nil {
			httpx.SendError(w, errors.New("username atau password salah"), http.StatusBadRequest)
			return
		}
		school, err := repositories.GetSchoolByStudentId(student.Id)
		if err != nil {
			httpx.SendError(w, errors.New("siswa tidak memiliki akses ke sekolah"), http.StatusBadRequest)
			return
		}
		claim := auth.JWTClaim{
			UserId:   student.Id,
			Username: student.Username,
			Email:    student.Email.String,
			Role:     "student",
			SchoolId: school.Id,
		}
		token := auth.CreateToken(claim)
		res := models.LoginResponse{
			Token:    token,
			Username: student.Username,
			Id:       student.Id,
			Role:     "student",
			SchoolId: school.Id,
		}
		httpx.SendData(w, res, http.StatusOK)
		return
	default:
		httpx.SendError(w, errors.New("tipe login tidak valid"), http.StatusBadRequest)
		return
	}
}

// HandlePostRefresh godoc
//
//	@Summary		Refresh auth session
//	@Description	Refresh aush session data
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		schemas.LoginRequest	true	"Login request body"
//	@Success		200		{object}	httpx.DataJsonResponse{data=schemas.LoginResponse}
//	@Failure		400		{object}	httpx.ErrorJsonResponse
//	@Failure		500		{object}	httpx.ErrorJsonResponse
//	@Router			/auth/login [post]
func HandlePostRefresh(w http.ResponseWriter, r *http.Request) {
	userIdClaim, _ := auth.GetJwtClaim(r, "user_id")
	userId := userIdClaim.(string)

	roleClaim, _ := auth.GetJwtClaim(r, "role")
	role := roleClaim.(string)

	schoolIdClaim, _ := auth.GetJwtClaim(r, "school_id")
	schoolId := schoolIdClaim.(string)

	switch role {
	case "admin":
		admin, err := user.GetAdminById(userId)
		if err != nil {
			if err == sql.ErrNoRows {
				httpx.SendError(w, errors.New("admin tidak ditemukan"), http.StatusNotFound)
				return
			}
			httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
			return
		}

		claim := auth.JWTClaim{
			UserId:   admin.Id,
			Username: admin.Username,
			Role:     "admin",
			Email:    admin.Email.String,
			SchoolId: schoolId,
		}
		token := auth.CreateToken(claim)
		res := schemas.AuthRefreshResponse{
			Token:    token,
			Username: admin.Username,
			Role:     "admin",
			Id:       admin.Id,
			SchoolId: schoolId,
		}
		httpx.SendData(w, res, http.StatusCreated)
		return
	case "student":
		student, err := user.GetStudentById(userId)
		if err != nil {
			if err == sql.ErrNoRows {
				httpx.SendError(w, errors.New("student tidak ditemukan"), http.StatusNotFound)
				return
			}
			httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
			return
		}

		claim := auth.JWTClaim{
			UserId:   student.Id,
			Username: student.Username,
			Role:     "student",
			Email:    student.Email.String,
			SchoolId: schoolId,
		}
		token := auth.CreateToken(claim)
		res := schemas.AuthRefreshResponse{
			Token:    token,
			Username: student.Username,
			Role:     "student",
			Id:       student.Id,
			SchoolId: schoolId,
		}
		httpx.SendData(w, res, http.StatusCreated)
		return
	}
}
