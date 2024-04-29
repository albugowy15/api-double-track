package repositories

import (
	"database/sql"
	"log"

	"github.com/albugowy15/api-double-track/db"
	"github.com/albugowy15/api-double-track/internal/models"
)

func GetStudentExpectations(studentId string) ([]models.ExpectationToAlternative, error) {
	var expectations []models.ExpectationToAlternative
	err := db.AppDB.Select(
		&expectations,
		`SELECT eta.id, eta.rank, eta.alternative_id, eta.expectation_id,
    a.alternative as alternative_name FROM expectations e
    INNER JOIN expectations_to_alternatives eta ON e.id = eta.expectation_id
    INNER JOIN alternatives a ON a.id = eta.alternative_id
    WHERE e.student_id = $1`, studentId)
	if err != nil {
		log.Println("db err: ", err)
	}
	return expectations, err
}

func CheckStudentExpectationExist(studentId string) bool {
	var expectationId string
	err := db.AppDB.Get(&expectationId, `SELECT id FROM expectations WHERE student_id = $1`, studentId)
	return err != sql.ErrNoRows
}

func SaveStudentExpectations(expectations []models.ExpectationData, studentId string) error {
	tx, err := db.AppDB.Beginx()
	if err != nil {
		return err
	}
	var expectationId string
	err = tx.QueryRowx(`INSERT INTO expectations (student_id) VALUES ($1) RETURNING id`, studentId).Scan(&expectationId)
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, expect := range expectations {
		_, err := tx.Exec(
			`INSERT INTO expectations_to_alternatives (expectation_id, alternative_id, rank) VALUES ($1, $2, $3)`,
			expectationId, expect.AlternativeId, expect.Rank,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func UpdateStudentExpectation(expectations []models.ExpectationData, studentId string) error {
	var expectationId string
	err := db.AppDB.Get(&expectationId, "SELECT id FROM expectations WHERE student_id = $1", studentId)
	if err != nil {
		return err
	}

	tx, err := db.AppDB.Beginx()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`DELETE FROM expectations_to_alternatives WHERE expectation_id = $1`, expectationId)
	if err != nil {
		tx.Rollback()
		return err
	}
	for _, expect := range expectations {
		_, err := tx.Exec(
			`INSERT INTO expectations_to_alternatives (expectation_id, alternative_id, rank) VALUES ($1, $2, $3)`,
			expectationId, expect.AlternativeId, expect.Rank,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func DeleteStudentExpectation(studentId string) error {
	var expectationId string
	err := db.AppDB.Get(&expectationId, "SELECT id FROM expectations WHERE student_id = $1", studentId)
	if err != nil {
		return err
	}

	tx, err := db.AppDB.Beginx()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`DELETE FROM expectations_to_alternatives WHERE expectation_id = $1`, expectationId)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(`DELETE FROM expectations WHERE id = $1`, expectationId)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}
