package repositories

import (
	"log"

	"github.com/albugowy15/api-double-track/db"
	"github.com/albugowy15/api-double-track/internal/models"
)

type QuestionRepository struct{}

var questionRepository *QuestionRepository

func GetQuestionRepository() *QuestionRepository {
	if questionRepository == nil {
		questionRepository = &QuestionRepository{}
	}
	return questionRepository
}

func (r *QuestionRepository) GetQuestions() ([]models.Question, error) {
	questions := []models.Question{}
	err := db.AppDB.Select(&questions, "SELECT id, question, description, category, code, number FROM questions ORDER BY number")
	if err != nil {
		log.Printf("err repository get questions: %v", err)
	}

	return questions, err
}
