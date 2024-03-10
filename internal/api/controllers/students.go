package controllers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth"
	"github.con/albugowy15/api-double-track/internal/pkg/repositories/user"
	"github.con/albugowy15/api-double-track/internal/pkg/utils"
)

func GetStudents(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	role, _ := claims["role"].(string)
	if role != "admin" {
		utils.SendError(w, "anda bukan admin", http.StatusForbidden)
		return
	}

	schoolId, _ := claims["school_id"].(string)

	students, err := user.GetStudentRepository().GetStudentsBySchool(schoolId)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}
	utils.SendJson(w, students, http.StatusOK)
}

func GetStudent(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	schoolId, _ := claims["school_id"].(string)
	studentIdParam := chi.URLParam(r, "studentId")
	student, err := user.GetStudentRepository().GetStudentBySchoolId(schoolId, studentIdParam)
	if err != nil {
		utils.SendError(w, "data siswa tidak ditemukan", http.StatusNotFound)
		return
	}
	utils.SendJson(w, student, http.StatusOK)
}
