package user

import (
	"log"

	"github.com/albugowy15/api-double-track/internal/pkg/db"
	"github.com/albugowy15/api-double-track/internal/pkg/models/user"
)

type StudentRepository struct{}

var studentRepository *StudentRepository

func GetStudentRepository() *StudentRepository {
	if studentRepository == nil {
		studentRepository = &StudentRepository{}
	}
	return studentRepository
}

func (r *StudentRepository) GetStudentByUsername(username string) (user.Student, error) {
	student := user.Student{}
	err := db.GetDb().Get(&student, "SELECT id, username, email, password, phone_number, fullname, nisn  FROM students WHERE username = $1", username)
	return student, err
}

func (r *StudentRepository) GetStudentsBySchool(schoolId string) ([]user.Student, error) {
	students := []user.Student{}
	err := db.GetDb().Select(&students, "SELECT id, username, email, phone_number, fullname, nisn FROM students WHERE school_id = $1", schoolId)
	return students, err
}

func (r *StudentRepository) GetStudentById(studentId string) (user.Student, error) {
	student := user.Student{}
	err := db.GetDb().Get(&student, "SELECT id, username, email, phone_number, fullname, nisn FROM students WHERE id = $1", studentId)
	return student, err
}

func (r *StudentRepository) GetStudentBySchoolId(schoolId string, studentId string) (user.Student, error) {
	student := user.Student{}
	err := db.GetDb().Get(&student, "SELECT id, username, email, phone_number, fullname, nisn FROM students WHERE school_id = $1 AND id = $2", schoolId, studentId)
	return student, err
}

func (r *StudentRepository) AddStudent(schoolId string, data user.Student) error {
	_, err := db.GetDb().Exec("INSERT INTO students (username, password, fullname, nisn, school_id) VALUES ($1, $2, $3, $4, $5)", data.Username, data.Password, data.Fullname, data.Nisn, schoolId)
	return err
}

func (r *StudentRepository) DeleteStudent(studentId string) error {
	_, err := db.GetDb().Exec("DELETE FROM students WHERE id = $1", studentId)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (r *StudentRepository) UpdateStudent(studentId string, data user.Student) error {
	_, err := db.GetDb().Exec(
		"UPDATE students SET fullname = $1, username = $2, nisn = $3, email = $4, phone_number = $5 WHERE id = $6",
		data.Fullname,
		data.Username,
		data.Nisn,
		data.Email,
		data.PhoneNumber,
		studentId,
	)
	if err != nil {
		log.Println(err)
	}
	return err
}
