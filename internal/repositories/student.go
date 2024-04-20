package repositories

import (
	"database/sql"
	"log"

	"github.com/albugowy15/api-double-track/db"
	"github.com/albugowy15/api-double-track/internal/models"
)

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
	tx, err := db.AppDB.Beginx()
	if err != nil {
		log.Println(err)
		return err
	}

	var ahpId int64
	err = tx.Get(&ahpId, `SELECT id FROM ahp WHERE student_id = $1`, studentId)
	if err == sql.ErrNoRows {
		// it means the student is not completed the questionnare yet
		// proceed to delete its row from students
		_, err = tx.Exec("DELETE FROM students WHERE id = $1", studentId)
		if err != nil {
			log.Println(err)
			tx.Rollback()
			return err
		}
		tx.Commit()
		return nil
	}

	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	// delete record from answers
	_, err = tx.Exec(`DELETE FROM answers WHERE student_id = $1`, studentId)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	// delete record from ahp_to_alternatives
	_, err = tx.Exec(`DELETE FROM ahp_to_alternatives WHERE ahp_id = $1`, ahpId)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	// TODO: delete record from topsis_to_alternatives

	// delete record from ahp
	_, err = tx.Exec(`DELETE FROM ahp WHERE student_id = $1`, studentId)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}
	// TODO: delete record from topsis

	// finally delete student record
	_, err = tx.Exec("DELETE FROM students WHERE id = $1", studentId)
	if err != nil {
		log.Println(err)
		tx.Rollback()
		return err
	}

	tx.Commit()
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
