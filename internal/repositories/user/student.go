package user

import (
	"log"

	"github.com/albugowy15/api-double-track/db"
	"github.com/albugowy15/api-double-track/internal/models/user"
)

func GetStudentByUsername(username string) (user.Student, error) {
	student := user.Student{}
	err := db.AppDB.Get(
		&student,
		"SELECT id, username, email, password, phone_number, fullname, nisn  FROM students WHERE username = $1",
		username,
	)
	if err != nil {
		log.Println(err)
	}
	return student, err
}

func GetStudentsBySchool(schoolId string) ([]user.Student, error) {
	students := []user.Student{}
	err := db.AppDB.Select(
		&students,
		"SELECT id, username, email, phone_number, fullname, nisn FROM students WHERE school_id = $1",
		schoolId,
	)
	if err != nil {
		log.Println(err)
	}
	return students, err
}

func GetStudentById(studentId string) (user.Student, error) {
	student := user.Student{}
	err := db.AppDB.Get(
		&student,
		"SELECT id, username, email, phone_number, fullname, nisn FROM students WHERE id = $1",
		studentId,
	)
	if err != nil {
		log.Println(err)
	}
	return student, err
}

func GetStudentBySchoolId(schoolId string, studentId string) (user.Student, error) {
	student := user.Student{}
	err := db.AppDB.Get(
		&student,
		"SELECT id, username, email, phone_number, fullname, nisn FROM students WHERE school_id = $1 AND id = $2",
		schoolId,
		studentId,
	)
	if err != nil {
		log.Println(err)
	}
	return student, err
}

func AddStudent(schoolId string, data user.Student) error {
	_, err := db.AppDB.Exec(
		`INSERT INTO students (username, password, fullname, nisn, school_id) VALUES ($1, $2, $3, $4, $5)`,
		data.Username,
		data.Password,
		data.Fullname,
		data.Nisn,
		schoolId,
	)
	if err != nil {
		log.Println(err)
	}
	return err
}

func DeleteStudent(studentId string) error {
	_, err := db.AppDB.Exec("DELETE FROM students WHERE id = $1", studentId)
	if err != nil {
		log.Println(err)
	}
	return err
}

func UpdateStudent(studentId string, data user.Student) error {
	_, err := db.AppDB.Exec(
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

func GetTotalStudents(schoolId string) (int64, error) {
	var totalStudent int64
	err := db.AppDB.Get(
		&totalStudent,
		`SELECT COUNT(*) as total_student FROM students WHERE school_id = $1`,
		schoolId,
	)
	if err != nil {
		log.Println(err)
	}
	return totalStudent, err
}
