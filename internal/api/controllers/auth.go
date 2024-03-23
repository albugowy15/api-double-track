package controllers

import (
	"errors"
	"log"
	"net/http"

	"github.com/albugowy15/api-double-track/internal/pkg/models"
	"github.com/albugowy15/api-double-track/internal/pkg/repositories"
	"github.com/albugowy15/api-double-track/internal/pkg/repositories/user"
	"github.com/albugowy15/api-double-track/internal/pkg/utils/httputil"
	"github.com/albugowy15/api-double-track/internal/pkg/utils/jwt"
	"github.com/albugowy15/api-double-track/internal/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

// Login godoc
//
//	@Summary		Login authentication
//	@Description	Login authentication for student and admin
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			body	body		schemas.LoginRequest	true	"Login request body"
//	@Success		200		{object}	httputil.DataJsonResponse{data=schemas.LoginResponse}
//	@Failure		400		{object}	httputil.ErrorJsonResponse
//	@Failure		500		{object}	httputil.ErrorJsonResponse
//	@Router			/auth/login [post]
func Login(w http.ResponseWriter, r *http.Request) {
	var body models.LoginRequest
	httputil.GetBody(w, r, &body)
	err := validator.ValidateLoginRequest(body)
	if err != nil {
		httputil.SendError(w, err, http.StatusBadRequest)
		return
	}

	switch body.Type {
	case "admin":
		s := user.GetAdminRepository()
		admin, err := s.GetAdminByUsername(body.Username)
		if err != nil {
			httputil.SendError(w, errors.New("username atau password salah"), http.StatusBadRequest)
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(body.Password)); err != nil {
			httputil.SendError(w, errors.New("username atau password salah"), http.StatusBadRequest)
			return
		}
		school, err := repositories.GetSchoolRepository().GetSchoolByAdminId(admin.Id)
		if err != nil {
			httputil.SendError(w, errors.New("admin tidak memiliki akses ke sekolah"), http.StatusBadRequest)
			return
		}
		claim := jwt.JWTClaim{
			UserId:   admin.Id,
			Username: admin.Username,
			Email:    admin.Email.String,
			Role:     "admin",
			SchoolId: school.Id,
		}
		token := jwt.CreateToken(claim)
		res := models.LoginResponse{
			Token:    token,
			Username: admin.Username,
			Id:       admin.Id,
			Role:     "admin",
			SchoolId: school.Id,
		}
		httputil.SendData(w, res, http.StatusOK)
		return
	case "student":
		s := user.GetStudentRepository()
		student, err := s.GetStudentByUsername(body.Username)
		if err != nil {
			log.Println(err)
			httputil.SendError(w, errors.New("username atau password salah"), http.StatusBadRequest)
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(body.Password)); err != nil {
			httputil.SendError(w, errors.New("username atau password salah"), http.StatusBadRequest)
			return
		}
		school, err := repositories.GetSchoolRepository().GetSchoolByStudentId(student.Id)
		if err != nil {
			httputil.SendError(w, errors.New("siswa tidak memiliki akses ke sekolah"), http.StatusBadRequest)
			return
		}
		claim := jwt.JWTClaim{
			UserId:   student.Id,
			Username: student.Username,
			Email:    student.Email.String,
			Role:     "student",
			SchoolId: student.Id,
		}
		token := jwt.CreateToken(claim)
		res := models.LoginResponse{
			Token:    token,
			Username: student.Username,
			Id:       student.Id,
			Role:     "student",
			SchoolId: school.Id,
		}
		httputil.SendData(w, res, http.StatusOK)
		return
	default:
		httputil.SendError(w, errors.New("tipe login tidak valid"), http.StatusBadRequest)
		return
	}
}
