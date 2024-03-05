package controllers

import (
	"net/http"

	"github.con/albugowy15/api-double-track/internal/pkg/models"
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
			utils.SendError(w, "user tidak ditemukan", http.StatusBadRequest)
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(body.Password)); err != nil {
			utils.SendError(w, "password salah", http.StatusBadRequest)
			return
		}
		token := jwt.CreateToken(jwt.JWTClaim{UserId: admin.Id, Username: admin.Username, Email: admin.Email, Role: "admin"})
		res := models.LoginResponse{
			Token: token,
		}
		utils.SendJson(w, res, http.StatusOK)
		return
	case "student":
		s := user.GetStudentRepository()
		student, err := s.GetStudentByUsername(body.Username)
		if err != nil {
			utils.SendError(w, "user tidak ditemukan", http.StatusBadRequest)
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(student.Password), []byte(body.Password)); err != nil {
			utils.SendError(w, "password salah", http.StatusBadRequest)
			return
		}
		token := jwt.CreateToken(jwt.JWTClaim{UserId: student.Id, Username: student.Username, Email: student.Email, Role: "student"})
		res := models.LoginResponse{
			Token: token,
		}
		utils.SendJson(w, res, http.StatusOK)
		return
	default:
		utils.SendError(w, "tipe login tidak valid", http.StatusBadRequest)
		return
	}
}
