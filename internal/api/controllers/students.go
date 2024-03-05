package controllers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.con/albugowy15/api-double-track/internal/pkg/repositories/user"
	"github.con/albugowy15/api-double-track/internal/pkg/utils"
)

func GetStudents(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	role, ok := claims["role"].(string)
	if !ok {
		utils.SendError(w, "token invalid", http.StatusBadRequest)
		return
	}
	log.Printf("role: %s", role)
	if role != "admin" {
		utils.SendError(w, "anda bukan admin", http.StatusForbidden)
		return
	}
	students, err := user.GetStudentRepository().GetStudents()
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.SendJson(w, students, http.StatusOK)
}

func GetStudent(w http.ResponseWriter, r *http.Request) {
	studentIdParam := chi.URLParam(r, "studentId")
	student, err := user.GetStudentRepository().GetStudentById(studentIdParam)
	if err != nil {
		log.Println(err)
		utils.SendError(w, "data siswa tidak ditemukan", http.StatusNotFound)
		return
	}
	utils.SendJson(w, student, http.StatusOK)
}
