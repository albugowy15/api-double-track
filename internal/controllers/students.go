package controllers

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/albugowy15/api-double-track/internal/models"
	"github.com/albugowy15/api-double-track/internal/repositories"
	"github.com/albugowy15/api-double-track/internal/validator"
	"github.com/albugowy15/api-double-track/pkg/auth"
	"github.com/albugowy15/api-double-track/pkg/crypto"
	"github.com/albugowy15/api-double-track/pkg/httpx"
	"github.com/go-chi/chi/v5"
	"github.com/lib/pq"
)

// HandleStudents godoc
//
//	@Summary		Get students
//	@Description	Get all students from a school
//	@Tags			Students
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		200				{object}	httpx.DataJsonResponse{data=[]schemas.Student}
//	@Failure		400				{object}	httpx.ErrorJsonResponse
//	@Failure		500				{object}	httpx.ErrorJsonResponse
//	@Router			/students [get]
func HandleGetStudents(w http.ResponseWriter, r *http.Request) {
	schoolIdClaim, _ := auth.GetJwtClaim(r, "school_id")
	schoolId := schoolIdClaim.(string)

	students, err := repositories.GetStudentsBySchool(schoolId)
	if err != nil {
		httpx.SendError(w, err, http.StatusBadRequest)
		return
	}

	httpx.SendData(w, students, http.StatusOK)
}

// HandleGetStudent godoc
//
//	@Summary		Get a student
//	@Description	Get a student from a school
//	@Tags			Students
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			studentId		path		string	true	"Id student"
//	@Success		200				{object}	httpx.DataJsonResponse{data=schemas.Student}
//	@Failure		400				{object}	httpx.ErrorJsonResponse
//	@Failure		404				{object}	httpx.ErrorJsonResponse
//	@Failure		500				{object}	httpx.ErrorJsonResponse
//	@Router			/students/{studentId} [get]
func HandleGetStudent(w http.ResponseWriter, r *http.Request) {
	schoolIdClaim, _ := auth.GetJwtClaim(r, "school_id")
	schoolId := schoolIdClaim.(string)
	studentIdParam := chi.URLParam(r, "studentId")
	student, err := repositories.GetStudentBySchoolId(schoolId, studentIdParam)
	if err != nil {
		if err == sql.ErrNoRows {
			httpx.SendError(w, errors.New("data siswa tidak ditemukan"), http.StatusNotFound)
			return
		}
		log.Println(err)
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	httpx.SendData(w, student, http.StatusOK)
}

// HandlePostStudent godoc
//
//	@Summary		Add a student
//	@Description	Add a student from a school
//	@Tags			Students
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string						true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			body			body		schemas.AddStudentRequest	true	"Add student request body"
//	@Success		201				{object}	httpx.MessageJsonResponse
//	@Failure		400				{object}	httpx.ErrorJsonResponse
//	@Failure		500				{object}	httpx.ErrorJsonResponse
//	@Router			/students [post]
func HandlePostStudent(w http.ResponseWriter, r *http.Request) {
	schoolIdClaim, _ := auth.GetJwtClaim(r, "school_id")
	schoolId := schoolIdClaim.(string)

	var body models.Student
	if err := httpx.GetBody(r, &body); err != nil {
		httpx.SendError(w, err, http.StatusBadRequest)
	}

	sanitizedBody, err := validator.ValidateAddStudent(body)
	if err != nil {
		httpx.SendError(w, err, http.StatusBadRequest)
		return
	}

	_, err = repositories.GetSchoolById(schoolId)
	if err != nil {
		if err == sql.ErrNoRows {
			httpx.SendError(w, errors.New("id sekolah tidak ditemukan"), http.StatusBadRequest)
			return
		}
		log.Println(err)
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	sanitizedBody.Username = sanitizedBody.Nisn
	password, err := crypto.HashStr(sanitizedBody.Nisn)
	if err != nil {
		httpx.SendError(w, err, http.StatusBadRequest)
		return
	}
	sanitizedBody.Password = password

	if err := repositories.AddStudent(schoolId, sanitizedBody); err != nil {
		if err, _ := err.(*pq.Error); err.Code.Class() == "23" {
			httpx.SendError(w, errors.New("nisn sudah terdaftar"), http.StatusBadRequest)
			return
		}
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	httpx.SendMessage(w, "berhasil menambah siswa", http.StatusCreated)
}

// HandleDeleteStudent godoc
//
//	@Summary		Delete a student
//	@Description	Delete a student from a school
//	@Tags			Students
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			body			body		schemas.DeleteStudentRequest	true	"Delete student request body"
//	@Success		201				{object}	httpx.MessageJsonResponse
//	@Failure		400				{object}	httpx.ErrorJsonResponse
//	@Failure		500				{object}	httpx.ErrorJsonResponse
//	@Router			/students [delete]
func HandleDeleteStudent(w http.ResponseWriter, r *http.Request) {
	var body models.Student
	if err := httpx.GetBody(r, &body); err != nil {
		httpx.SendError(w, err, http.StatusBadRequest)
	}

	if len(body.Id) == 0 {
		httpx.SendError(w, errors.New("id siswa wajib diisi"), http.StatusBadRequest)
		return
	}

	// check id exist
	student, err := repositories.GetStudentById(body.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			httpx.SendError(w, errors.New("siswa tidak ditemukan"), http.StatusBadRequest)
			return
		}
		log.Println(err)
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	studentSchool, err := repositories.GetSchoolByStudentId(student.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			httpx.SendError(w, errors.New("sekolah tidak ditemukan"), http.StatusBadRequest)
			return
		}
		log.Println(err)
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	// check student school id same with admin school id
	schoolIdClaim, _ := auth.GetJwtClaim(r, "school_id")
	schoolId := schoolIdClaim.(string)
	if schoolId != studentSchool.Id {
		httpx.SendError(w, errors.New("tidak memiliki akses ke sekolah"), http.StatusUnauthorized)
		return
	}

	err = repositories.DeleteStudent(student.Id)
	if err != nil {
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	httpx.SendMessage(w, "berhasil menghapus data siswa", http.StatusOK)
}

// HandlePatchStudent godoc
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
//	@Success		201				{object}	httpx.MessageJsonResponse
//	@Failure		400				{object}	httpx.ErrorJsonResponse
//	@Failure		500				{object}	httpx.ErrorJsonResponse
//	@Router			/students/{studentId} [patch]
func HandlePatchStudent(w http.ResponseWriter, r *http.Request) {
	studentIdParam := chi.URLParam(r, "studentId")

	var body models.Student
	if err := httpx.GetBody(r, &body); err != nil {
		httpx.SendError(w, err, http.StatusBadRequest)
	}
	sanitizedBody, err := validator.ValidateUpdateStudent(body)
	if err != nil {
		httpx.SendError(w, err, http.StatusBadRequest)
		return
	}

	// check id exist
	student, err := repositories.GetStudentById(studentIdParam)
	if err != nil {
		if err == sql.ErrNoRows {
			httpx.SendError(w, errors.New("siswa tidak ditemukan"), http.StatusBadRequest)
			return
		}
		log.Println(err)
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	// check same school id
	studentSchool, err := repositories.GetSchoolByStudentId(student.Id)
	if err != nil {
		if err == sql.ErrNoRows {
			httpx.SendError(w, errors.New("sekolah tidak ditemukan"), http.StatusBadRequest)
			return
		}
		log.Println(err)
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	schoolIdClaim, _ := auth.GetJwtClaim(r, "school_id")
	schoolId := schoolIdClaim.(string)
	if schoolId != studentSchool.Id {
		httpx.SendError(w, errors.New("tidak memiliki akses ke sekolah"), http.StatusUnauthorized)
		return
	}

	// save to db
	err = repositories.UpdateStudent(student.Id, sanitizedBody)
	if err != nil {
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	// success
	httpx.SendMessage(w, "berhasil memperbarui data siswa", http.StatusOK)
}

// HandleGetStudentProfile godoc
//
//	@Summary		Get a student profile
//	@Description	Get a student profile
//	@Tags			Students
//	@Tags			Student
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		200				{object}	httpx.DataJsonResponse{data=schemas.StudentProfile}
//	@Failure		400				{object}	httpx.ErrorJsonResponse
//	@Failure		500				{object}	httpx.ErrorJsonResponse
//	@Router			/students/profile [get]
func HandleGetStudentProfile(w http.ResponseWriter, r *http.Request) {
	// get student id from token
	studentIdClaim, _ := auth.GetJwtClaim(r, "user_id")
	studentId := studentIdClaim.(string)

	// Get student from db by student id
	student, err := repositories.GetStudentById(studentId)
	if err != nil {
		if err == sql.ErrNoRows {
			httpx.SendError(w, errors.New("data siswa tidak ditemukan"), http.StatusNotFound)
			return
		}
		log.Println(err)
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	school, err := repositories.GetSchoolByStudentId(studentId)
	if err != nil {
		if err == sql.ErrNoRows {
			httpx.SendError(w, errors.New("data sekolah tidak ditemukan"), http.StatusNotFound)
			return
		}
		log.Println(err)
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	profile := models.StudentProfile{
		Id:          student.Id,
		Username:    student.Username,
		Fullname:    student.Fullname,
		Email:       student.Email,
		PhoneNumber: student.PhoneNumber,
		Nisn:        student.Nisn,
		School:      school.Name,
	}
	httpx.SendData(w, profile, http.StatusOK)
}

// HandlePatchStudentProfile godoc
//
//	@Summary		Update a student profile
//	@Description	Update a student profile
//	@Tags			Students
//	@Tags			Student
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string							true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			body			body		schemas.UpdateStudentRequest	true	"Update student profile body request"
//	@Success		201				{object}	httpx.MessageJsonResponse
//	@Failure		400				{object}	httpx.ErrorJsonResponse
//	@Failure		500				{object}	httpx.ErrorJsonResponse
//	@Router			/students/profile [patch]
func HandlePatchStudentProfile(w http.ResponseWriter, r *http.Request) {
	studentIdClaim, _ := auth.GetJwtClaim(r, "user_id")
	studentId := studentIdClaim.(string)

	var body models.Student
	if err := httpx.GetBody(r, &body); err != nil {
		httpx.SendError(w, err, http.StatusBadRequest)
	}
	sanitizedBody, err := validator.ValidateUpdateStudent(body)
	if err != nil {
		httpx.SendError(w, err, http.StatusBadRequest)
		return
	}

	// check id exist
	student, err := repositories.GetStudentById(studentId)
	if err != nil {
		if err == sql.ErrNoRows {
			httpx.SendError(w, errors.New("siswa tidak ditemukan"), http.StatusBadRequest)
			return
		}
		log.Println(err)
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	// save to db
	err = repositories.UpdateStudent(student.Id, sanitizedBody)
	if err != nil {
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	// success
	httpx.SendMessage(w, "berhasil memperbarui profil", http.StatusOK)
}
