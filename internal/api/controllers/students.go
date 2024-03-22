package controllers

import (
	"log"
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

// GetStudents godoc
//
//	@Summary		Get students
//	@Description	Get all students from a school
//	@Tags			Students
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		200				{object}	utils.DataJsonResponse{data=[]schemas.Student}
//	@Failure		400				{object}	utils.ErrorJsonResponse
//	@Failure		500				{object}	utils.ErrorJsonResponse
//	@Router			/students [get]
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

// GetStudent godoc
//
//	@Summary		Get a student
//	@Description	Get a student from a school
//	@Tags			Students
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			studentId		path		string	true	"Id student"
//	@Success		200				{object}	utils.DataJsonResponse{data=schemas.Student}
//	@Failure		400				{object}	utils.ErrorJsonResponse
//	@Failure		500				{object}	utils.ErrorJsonResponse
//	@Router			/students/{studentId} [get]
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

// AddStudent godoc
//
//	@Summary		Add a student
//	@Description	Add a student from a school
//	@Tags			Students
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string						true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			body			body		schemas.AddStudentRequest	true	"Add student request body"
//	@Success		201				{object}	schemas.MessageResponse
//	@Failure		400				{object}	utils.ErrorJsonResponse
//	@Failure		500				{object}	utils.ErrorJsonResponse
//	@Router			/students [post]
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

// DeleteStudent godoc
//
//	@Summary		Delete a student
//	@Description	Delete a student from a school
//	@Tags			Students
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			body			body		schemas.DeleteStudentRequest	true	"Delete student request body"
//	@Success		201				{object}	schemas.MessageResponse
//	@Failure		400				{object}	utils.ErrorJsonResponse
//	@Failure		500				{object}	utils.ErrorJsonResponse
//	@Router			/students [delete]
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

// UpdateStudent godoc
//
//	@Summary		Update a student
//	@Description	Update a student from a school
//	@Tags			Students
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			studentId		path		string							true	"Update student id"
//	@Param			body			body		schemas.UpdateStudentRequest	true	"Update student request body"
//	@Success		201				{object}	schemas.MessageResponse
//	@Failure		400				{object}	utils.ErrorJsonResponse
//	@Failure		500				{object}	utils.ErrorJsonResponse
//	@Router			/students/{studentId} [patch]
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

// GetStudentProfile godoc
//
//	@Summary		Get a student profile
//	@Description	Get a student profile
//	@Tags			Students
//	@Tags			Student
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		200				{object}	utils.DataJsonResponse{data=schemas.StudentProfile}
//	@Failure		400				{object}	utils.ErrorJsonResponse
//	@Failure		500				{object}	utils.ErrorJsonResponse
//	@Router			/students/profile [get]
func GetStudentProfile(w http.ResponseWriter, r *http.Request) {
	// get student id from token
	studentIdClaim, _ := jwt.GetJwtClaim(r, "user_id")
	studentId := studentIdClaim.(string)

	// Get student from db by student id
	s := user.GetStudentRepository()
	student, err := s.GetStudentById(studentId)
	if err != nil {
		log.Printf("err get student profile: %v", err)
		utils.SendError(w, "data siswa tidak ditemukan", http.StatusNotFound)
		return
	}
	school, err := repositories.GetSchoolRepository().GetSchoolByStudentId(studentId)
	if err != nil {
		log.Printf("err get student school: %v", err)
		utils.SendError(w, "data sekolah tidak ditemukan", http.StatusNotFound)
		return
	}
	profile := userModel.StudentProfile{
		Id:          student.Id,
		Username:    student.Username,
		Fullname:    student.Fullname,
		Email:       student.Email,
		PhoneNumber: student.PhoneNumber,
		Nisn:        student.Nisn,
		School:      school.Name,
	}
	utils.SendJson(w, profile, http.StatusOK)
}

// UpdateStudentProfile godoc
//
//	@Summary		Update a student profile
//	@Description	Update a student profile
//	@Tags			Students
//	@Tags			Student
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			body			body		schemas.UpdateStudentRequest	true	"Update student profile body request"
//	@Success		201				{object}	schemas.MessageResponse
//	@Failure		400				{object}	utils.ErrorJsonResponse
//	@Failure		500				{object}	utils.ErrorJsonResponse
//	@Router			/students/profile [patch]
func UpdateStudentProfile(w http.ResponseWriter, r *http.Request) {
	studentIdClaim, _ := jwt.GetJwtClaim(r, "user_id")
	studentId := studentIdClaim.(string)

	var body userModel.Student
	utils.GetBody(w, r, &body)
	sanitizedBody, err := validator.ValidateUpdateStudent(body)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// check id exist
	s := user.GetStudentRepository()
	student, err := s.GetStudentById(studentId)
	if err != nil {
		utils.SendError(w, "siswa tidak ditemukan", http.StatusBadRequest)
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
		Message: "berhasil memperbarui profil",
	}
	utils.SendJson(w, res, http.StatusOK)
}
