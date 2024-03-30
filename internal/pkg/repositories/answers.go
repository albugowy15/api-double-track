package repositories

import (
	"log"

	"github.com/albugowy15/api-double-track/internal/pkg/db"
	"github.com/albugowy15/api-double-track/internal/pkg/models"
	"github.com/jmoiron/sqlx"
)

type AnswersRepository struct{}

var answersRepository *AnswersRepository

func GetAnswersRepository() *AnswersRepository {
	if answersRepository == nil {
		answersRepository = &AnswersRepository{}
	}
	return answersRepository
}

func (r *AnswersRepository) SaveAnswersTx(answers []models.Answer, tx *sqlx.Tx) error {
	_, err := tx.NamedExec(`INSERT INTO answers (student_id, question_id, answer) VALUES (:student_id, :question_id, :answer)`, answers)
	if err != nil {
		log.Println("db err:", err)
	}
	return err
}

func (r *AnswersRepository) GetAnswersByStudentId(studentId string) ([]models.Answer, error) {
	answers := []models.Answer{}
	err := db.GetDb().Select(&answers, "SELECT id, student_id, question_id, answer FROM answers WHERE student_id = $1", studentId)
	if err != nil {
		log.Println("db err:", err)
	}
	return answers, err
}
