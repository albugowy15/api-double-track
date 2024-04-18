package repositories

import (
	"log"

	"github.com/albugowy15/api-double-track/db"
	"github.com/albugowy15/api-double-track/internal/models"
)

func GetQuestions() ([]models.Question, error) {
	questions := []models.Question{}
	err := db.AppDB.Select(&questions, "SELECT id, question, description, category, code, number FROM questions ORDER BY number")
	if err != nil {
		log.Printf("err repository get questions: %v", err)
	}

	return questions, err
}
