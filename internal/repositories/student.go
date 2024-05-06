package repositories

import (
	"database/sql"
	"log"

	"github.com/albugowy15/api-double-track/db"
	"github.com/albugowy15/api-double-track/internal/models"
)

func IsUniqueStudentUsername(username string) (bool, error) {
	var studentId string
	err := db.AppDB.Get(&studentId, "SELECT id FROM students WHERE username = $1", username)
	if err == nil {
		return false, nil
	}
	if err == sql.ErrNoRows {
		return true, nil
	}
	log.Println("db err:", err)
	return false, err
}

func IsUniqueStudentUsernameFromId(studentId string, username string) (bool, error) {
	var id string
	err := db.AppDB.Get(&id, "SELECT id FROM students WHERE username = $1 AND id != $2", username, studentId)
	if err == nil {
		return false, nil
	}
	if err == sql.ErrNoRows {
		return true, nil
	}
	log.Println("db err:", err)
	return false, err
}

func IsUniqueStudentEmail(email string) (bool, error) {
	var studentId string
	err := db.AppDB.Get(&studentId, "SELECT id FROM students WHERE email = $1", email)
	if err == nil {
		return false, nil
	}
	if err == sql.ErrNoRows {
		return true, nil
	}
	log.Println("db err:", err)
	return false, err
}

func IsUniqueStudentEmailFromId(studentId string, email string) (bool, error) {
	var id string
	err := db.AppDB.Get(&id, "SELECT id FROM students WHERE email = $1 AND id != $2", email, studentId)
	if err == nil {
		return false, nil
	}
	if err == sql.ErrNoRows {
		return true, nil
	}
	log.Println("db err:", err)
	return false, err
}

func IsUniqueStudentPhoneNumber(phoneNumber string) (bool, error) {
	var studentId string
	err := db.AppDB.Get(&studentId, "SELECT id FROM students WHERE phone_number = $1", phoneNumber)
	if err == nil {
		return false, nil
	}
	if err == sql.ErrNoRows {
		return true, nil
	}
	log.Println("db err:", err)
	return false, err
}

func IsUniqueStudentPhoneNumberFromId(studentId string, phoneNumber string) (bool, error) {
	var id string
	err := db.AppDB.Get(&id, "SELECT id FROM students WHERE phone_number = $1 AND id != $2", phoneNumber, studentId)
	if err == nil {
		return false, nil
	}
	if err == sql.ErrNoRows {
		return true, nil
	}
	log.Println("db err:", err)
	return false, err
}

func IsUniqueStudentNisn(nisn string) (bool, error) {
	var studentId string
	err := db.AppDB.Get(&studentId, "SELECT id FROM students WHERE nisn = $1", nisn)
	if err == nil {
		return false, nil
	}
	if err == sql.ErrNoRows {
		return true, nil
	}
	log.Println("db err:", err)
	return false, err
}

func IsUniqueStudentNisnFromId(studentId string, nisn string) (bool, error) {
	var id string
	err := db.AppDB.Get(&id, "SELECT id FROM students WHERE nisn = $1 AND id != $2", nisn, studentId)
	if err == nil {
		return false, nil
	}
	if err == sql.ErrNoRows {
		return true, nil
	}
	log.Println("db err:", err)
	return false, err
}

func GetStudentByUsername(username string) (models.Student, error) {
	student := models.Student{}
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

func GetStudentsBySchool(schoolId string) ([]models.Student, error) {
	students := []models.Student{}
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

func GetStudentById(studentId string) (models.Student, error) {
	student := models.Student{}
	err := db.AppDB.Get(
		&student,
		"SELECT id, username, email, phone_number, password, fullname, nisn FROM students WHERE id = $1",
		studentId,
	)
	if err != nil {
		log.Println(err)
		return student, err
	}
	return student, nil
}

func GetStudentBySchoolId(schoolId string, studentId string) (models.Student, error) {
	student := models.Student{}
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

func AddStudent(schoolId string, data models.Student) error {
	_, err := db.AppDB.Exec(
		`INSERT INTO students (username, password, fullname, email, nisn, school_id) VALUES ($1, $2, $3, $4, $5, $6)`,
		data.Username,
		data.Password,
		data.Fullname,
		data.Email,
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
		return err
	}

	return nil
}

func UpdateStudent(studentId string, data models.Student) error {
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

func UpdateStudentPassword(studentId string, hashedPassword string) error {
	_, err := db.AppDB.Exec(`UPDATE students SET password = $1 WHERE id = $2`, hashedPassword, studentId)
	if err != nil {
		log.Println("db err: ", err)
		return err
	}
	return nil
}
