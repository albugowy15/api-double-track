package validator_test

import (
	"testing"

	"github.com/albugowy15/api-double-track/internal/pkg/models"
	"github.com/albugowy15/api-double-track/internal/pkg/validator"
)

func TestValidateAnswer(t *testing.T) {
	body := []models.SubmitAnswerRequest{}
	for i := 1; i <= 14; i++ {
		item := models.SubmitAnswerRequest{
			Id:     i,
			Number: i,
			Answer: "2",
		}
		body = append(body, item)
	}
	for i := 15; i <= 24; i++ {
		item := models.SubmitAnswerRequest{
			Id:     i,
			Number: i,
			Answer: "9",
		}
		body = append(body, item)
	}

	err := validator.ValidateSubmitAnswer(body)
	if err != nil {
		t.Errorf("expect no err, got: %v", err)
	}

	body[4].Answer = "7"
	err = validator.ValidateSubmitAnswer(body)
	if err != validator.ErrInvalidAnswer {
		t.Errorf("expect err to be: %v, got: %v", validator.ErrInvalidAnswer, err)
	}
	body[4].Answer = "2"

	body[17].Answer = "4"
	if err != validator.ErrInvalidAnswer {
		t.Errorf("expect err to be: %v, got: %v", validator.ErrInvalidAnswer, err)
	}

	body[3].Id = 345
	err = validator.ValidateSubmitAnswer(body)
	if err != validator.ErrInvalidId {
		t.Errorf("expect err to be: %v, got: %v", validator.ErrInvalidId, err)
	}
	body[3].Id = 4

	body[5].Id = 12
	err = validator.ValidateSubmitAnswer(body)
	if err != validator.ErrDuplicateId {
		t.Errorf("expect err to be: %v, got: %v", validator.ErrDuplicateId, err)
	}
	body[5].Id = 6

	body = append(body[:4], body[4+1:]...)
	err = validator.ValidateSubmitAnswer(body)
	if err != validator.ErrAnswerLen {
		t.Errorf("expect err to be: %v, got: %v", validator.ErrAnswerLen, err)
	}
}
