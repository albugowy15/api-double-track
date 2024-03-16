package repositories

import (
	"log"

	"github.com/albugowy15/api-double-track/internal/pkg/db"
	"github.com/albugowy15/api-double-track/internal/pkg/models"
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
	err := db.GetDb().Select(&questions, "SELECT id, question, description, category, code, number FROM questions ORDER BY number")
	if err != nil {
		log.Printf("err repository get questions: %v", err)
	}

	return questions, err
}
