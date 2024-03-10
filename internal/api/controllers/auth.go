package controllers

import (
	"net/http"

	"github.con/albugowy15/api-double-track/internal/pkg/models"
	"github.con/albugowy15/api-double-track/internal/pkg/repositories"
	"github.con/albugowy15/api-double-track/internal/pkg/repositories/user"
	"github.con/albugowy15/api-double-track/internal/pkg/utils"
	"github.con/albugowy15/api-double-track/internal/pkg/utils/jwt"
	"github.con/albugowy15/api-double-track/internal/pkg/validator"
	"golang.org/x/crypto/bcrypt"
)

func Login(w http.ResponseWriter, r *http.Request) {
	var body models.LoginRequest
	utils.GetBody(w, r, &body)
	err := validator.ValidateLoginRequest(body)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	switch body.Type {
	case "admin":
		s := user.GetAdminRepository()
		admin, err := s.GetAdminByUsername(body.Username)
		if err != nil {
			utils.SendError(w, "username atau password salah", http.StatusBadRequest)
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(body.Password)); err != nil {
			utils.SendError(w, "username atau password salah", http.StatusBadRequest)
			return
		}
		school, err := repositories.GetSchoolRepository().GetSchoolByAdminId(admin.Id)
		if err != nil {
			utils.SendError(w, err.Error(), http.StatusBadRequest)
			return
		}
		claim := jwt.JWTClaim{
			UserId:   admin.Id,
			Username: admin.Username,
			Email:    admin.Email,
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
		utils.SendJson(w, res, http.StatusOK)
		return
	case "student":
		s := user.GetStudentRepository()
		student, err := s.GetStudentByUsername(body.Username)
		if err != nil {
			utils.SendError(w, "username atau password salah", http.StatusBadRequest)
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(body.Password)); err != nil {
			utils.SendError(w, "username atau password salah", http.StatusBadRequest)
			return
		}
		school, err := repositories.GetSchoolRepository().GetSchoolByStudentId(student.Id)
		claim := jwt.JWTClaim{
			UserId:   student.Id,
			Username: student.Username,
			Email:    student.Email,
			Role:     "student",
			SchoolId: student.Id,
		}
		token := jwt.CreateToken(claim)
		if err != nil {
			utils.SendError(w, "siswa tidak memiliki akses ke sekolah", http.StatusBadRequest)
			return
		}
		res := models.LoginResponse{
			Token:    token,
			Username: student.Username,
			Id:       student.Id,
			Role:     "student",
			SchoolId: school.Id,
		}
		utils.SendJson(w, res, http.StatusOK)
		return
	default:
		utils.SendError(w, "tipe login tidak valid", http.StatusBadRequest)
		return
	}
}
