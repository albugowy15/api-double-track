package repositories

import (
	"log"

	"github.com/albugowy15/api-double-track/internal/pkg/db"
	"github.com/albugowy15/api-double-track/internal/pkg/models"
)

type AnswersRepository struct{}

var answersRepository *AnswersRepository

func GetAnswersRepository() *AnswersRepository {
	if answersRepository == nil {
		answersRepository = &AnswersRepository{}
	}
	return answersRepository
}

func (r *AnswersRepository) SaveAnswers(answers []models.Answer) error {
	tx, err := db.GetDb().Beginx()
	if err != nil {
		log.Fatal(err)
	}
	_, err = tx.NamedExec(`INSERT INTO answers (student_id, question_id, answer) VALUES (:student_id, :question_id, :answer)`, answers)
	if err != nil {
		tx.Rollback()
		return err
	}
	// for _, answer := range answers {
	// 	_, err := tx.Exec("INSERT INTO answers (student_id, question_id, answer) VALUES ($1, $2, $3)", answer.StudentId, answer.QuestionId, answer.Answer)
	// 	if err != nil {
	// 		tx.Rollback()
	// 		return err
	// 	}
	// }
	tx.Commit()
	return nil
}
