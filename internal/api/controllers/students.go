package controllers

import (
	"net/http"

	"github.com/albugowy15/api-double-track/internal/pkg/models"
	userModel "github.com/albugowy15/api-double-track/internal/pkg/models/user"
	"github.com/albugowy15/api-double-track/internal/pkg/repositories"
	"github.com/albugowy15/api-double-track/internal/pkg/repositories/user"
	"github.com/albugowy15/api-double-track/internal/pkg/utils"
	"github.com/albugowy15/api-double-track/internal/pkg/utils/jwt"
	"github.com/albugowy15/api-double-track/internal/pkg/validator"
	"github.com/go-chi/chi/v5"
	"github.com/lib/pq"
)

func GetStudents(w http.ResponseWriter, r *http.Request) {
	schoolIdClaim, _ := jwt.GetJwtClaim(r, "school_id")
	schoolId := schoolIdClaim.(string)

	students, err := user.GetStudentRepository().GetStudentsBySchool(schoolId)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	utils.SendJson(w, students, http.StatusOK)
}

func GetStudent(w http.ResponseWriter, r *http.Request) {
	schoolIdClaim, _ := jwt.GetJwtClaim(r, "school_id")
	schoolId := schoolIdClaim.(string)
	studentIdParam := chi.URLParam(r, "studentId")
	student, err := user.GetStudentRepository().GetStudentBySchoolId(schoolId, studentIdParam)
	if err != nil {
		utils.SendError(w, "data siswa tidak ditemukan", http.StatusNotFound)
		return
	}
	utils.SendJson(w, student, http.StatusOK)
}

func AddStudent(w http.ResponseWriter, r *http.Request) {
	schoolIdClaim, _ := jwt.GetJwtClaim(r, "school_id")
	schoolId := schoolIdClaim.(string)

	var body userModel.Student
	utils.GetBody(w, r, &body)

	sanitizedBody, err := validator.ValidateAddStudent(body)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = repositories.GetSchoolRepository().GetSchoolById(schoolId)
	if err != nil {
		utils.SendError(w, "id sekolah tidak ditemukan", http.StatusBadRequest)
		return
	}

	sanitizedBody.Username = sanitizedBody.Nisn
	password, err := utils.HashStr(sanitizedBody.Nisn)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusBadGateway)
		return
	}
	sanitizedBody.Password = password

	if err := user.GetStudentRepository().AddStudent(schoolId, sanitizedBody); err != nil {
		if err, _ := err.(*pq.Error); err.Code.Class() == "23" {
			utils.SendError(w, "nisn sudah terdaftar", http.StatusBadRequest)
			return
		}
		utils.SendError(w, "internal server error", http.StatusInternalServerError)
		return
	}
	res := models.MessageResponse{
		Message: "berhasil menambah siswa",
	}
	utils.SendJson(w, res, http.StatusCreated)
}

// check role
func DeleteStudent(w http.ResponseWriter, r *http.Request) {
	var body userModel.Student
	utils.GetBody(w, r, &body)

	if len(body.Id) == 0 {
		utils.SendError(w, "id siswa wajib diisi", http.StatusBadRequest)
		return
	}

	// check id exist
	s := user.GetStudentRepository()
	student, err := s.GetStudentById(body.Id)
	if err != nil {
		utils.SendError(w, "siswa tidak ditemukan", http.StatusBadRequest)
		return
	}

	studentSchool, err := repositories.GetSchoolRepository().GetSchoolByStudentId(student.Id)
	if err != nil {
		utils.SendError(w, "sekolah tidak ditemukan", http.StatusBadRequest)
		return
	}

	// check student school id same with admin school id
	schoolIdClaim, _ := jwt.GetJwtClaim(r, "school_id")
	schoolId := schoolIdClaim.(string)
	if schoolId != studentSchool.Id {
		utils.SendError(w, "tidak memiliki akses ke sekolah", http.StatusUnauthorized)
		return
	}

	err = s.DeleteStudent(student.Id)
	if err != nil {
		utils.SendError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	res := models.MessageResponse{
		Message: "berhasil menghapus data siswa",
	}
	utils.SendJson(w, res, http.StatusOK)
}

func UpdateStudent(w http.ResponseWriter, r *http.Request) {
	studentIdParam := chi.URLParam(r, "studentId")

	var body userModel.Student
	utils.GetBody(w, r, &body)
	sanitizedBody, err := validator.ValidateUpdateStudent(body)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check id exist
	s := user.GetStudentRepository()
	student, err := s.GetStudentById(studentIdParam)
	if err != nil {
		utils.SendError(w, "siswa tidak ditemukan", http.StatusBadRequest)
		return
	}

	// check same school id
	studentSchool, err := repositories.GetSchoolRepository().GetSchoolByStudentId(student.Id)
	if err != nil {
		utils.SendError(w, "sekolah tidak ditemukan", http.StatusBadRequest)
		return
	}
	schoolIdClaim, _ := jwt.GetJwtClaim(r, "school_id")
	schoolId := schoolIdClaim.(string)
	if schoolId != studentSchool.Id {
		utils.SendError(w, "tidak memiliki akses ke sekolah", http.StatusUnauthorized)
		return
	}

	// save to db
	err = s.UpdateStudent(student.Id, sanitizedBody)
	if err != nil {
		utils.SendError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// success
	res := models.MessageResponse{
		Message: "berhail memperbarui data siswa",
	}
	utils.SendJson(w, res, http.StatusOK)
}
