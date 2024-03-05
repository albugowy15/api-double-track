package user

import (
	"github.con/albugowy15/api-double-track/internal/pkg/db"
	"github.con/albugowy15/api-double-track/internal/pkg/models/user"
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

func (r *StudentRepository) GetStudents() ([]user.Student, error) {
	students := []user.Student{}
	err := db.GetDb().Select(&students, "SELECT id, username, email, phone_number, fullname, nisn FROM students")
	return students, err
}

func (r *StudentRepository) GetStudentById(studentId string) (user.Student, error) {
	student := user.Student{}
	err := db.GetDb().Get(&student, "SELECT id, username, email, phone_number, fullname, nisn FROM students WHERE id = $1", studentId)
	return student, err
}
