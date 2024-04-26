package repositories

import (
	"log"

	"github.com/albugowy15/api-double-track/db"
	"github.com/albugowy15/api-double-track/internal/models"
	"github.com/jmoiron/sqlx"
)

func SaveAnswersTx(answers []models.Answer, tx *sqlx.Tx) error {
	_, err := tx.NamedExec(`INSERT INTO answers (student_id, question_id, answer) VALUES (:student_id, :question_id, :answer)`, answers)
	if err != nil {
		log.Println("db err:", err)
	}
	return err
}

func GetAnswersByStudentId(studentId string) ([]models.Answer, error) {
	answers := []models.Answer{}
	err := db.AppDB.Select(&answers, "SELECT id, student_id, question_id, answer FROM answers WHERE student_id = $1", studentId)
	if err != nil {
		log.Println("db err:", err)
	}
	return answers, err
}

func DeleteAnswers(studentId string) error {
	tx, err := db.AppDB.Beginx()
	if err != nil {
		log.Println("db err:", err)
		return err
	}
	_, err = tx.Exec("DELETE FROM answers WHERE student_id = $1", studentId)
	if err != nil {
		tx.Rollback()
		log.Println("db err:", err)
		return err
	}
	var ahpId int32
	err = tx.QueryRowx("SELECT id FROM ahp WHERE student_id = $1", studentId).Scan(&ahpId)
	if err != nil {
		tx.Rollback()
		log.Println("db err:", err)
		return err
	}
	var topsisId int32
	err = tx.QueryRowx("SELECT id FROM topsis WHERE student_id = $1", studentId).Scan(&topsisId)
	if err != nil {
		tx.Rollback()
		log.Println("db err:", err)
		return err
	}
	_, err = tx.Exec("DELETE FROM ahp_to_alternatives WHERE ahp_id = $1", ahpId)
	if err != nil {
		tx.Rollback()
		log.Println("db err:", err)
		return err
	}
	_, err = tx.Exec("DELETE FROM topsis_to_alternatives WHERE topsis_id = $1", topsisId)
	if err != nil {
		tx.Rollback()
		log.Println("db err:", err)
		return err
	}
	_, err = tx.Exec("DELETE FROM ahp WHERE student_id = $1", studentId)
	if err != nil {
		tx.Rollback()
		log.Println("db err:", err)
		return err
	}
	_, err = tx.Exec("DELETE FROM topsis WHERE student_id = $1", studentId)
	if err != nil {
		tx.Rollback()
		log.Println("db err:", err)
		return err
	}
	tx.Commit()
	return nil
}
