package validator_test

import (
	"strconv"
	"testing"

	"github.com/albugowy15/api-double-track/internal/pkg/validator"
)

func TestValidateAnswer(t *testing.T) {
	body := map[string]string{}
	for i := 1; i <= 14; i++ {
		key := strconv.Itoa(i)
		body[key] = "2"
	}
	for i := 15; i <= 24; i++ {
		key := strconv.Itoa(i)
		body[key] = "9"
	}

	err := validator.ValidateSubmitAnswer(body)
	if err != nil {
		t.Errorf("expect no err, got: %v", err)
	}

	body["3"] = "7"
	err = validator.ValidateSubmitAnswer(body)
	if err != validator.ErrInvalidAnswer {
		t.Errorf("expect err to be: %v, got: %v", validator.ErrAnswerLen, err)
	}

	body["3"] = "2"
	body["15"] = "4"
	if err != validator.ErrInvalidAnswer {
		t.Errorf("expect err to be: %v, got: %v", validator.ErrAnswerLen, err)
	}

	delete(body, "4")
	err = validator.ValidateSubmitAnswer(body)
	if err != validator.ErrAnswerLen {
		t.Errorf("expect err to be: %v, got: %v", validator.ErrAnswerLen, err)
	}

	body["invalid"] = "2"
	err = validator.ValidateSubmitAnswer(body)
	if err != validator.ErrInvalidId {
		t.Errorf("expect err to be: %v, got: %v", validator.ErrInvalidId, err)
	}
}
