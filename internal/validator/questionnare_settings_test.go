package validator_test

import (
	"encoding/json"
	"testing"

	"github.com/albugowy15/api-double-track/internal/models"
	"github.com/albugowy15/api-double-track/internal/validator"
	"github.com/guregu/null/v5"
)

func TestValidateSettingOptions(t *testing.T) {
	settingItem := null.Int16From(1)

	err := validator.ValidateSettingItem(settingItem)
	if err != nil {
		t.Errorf("expect no err, got: %v", err)
	}

	settingItem.Valid = false
	err = validator.ValidateSettingItem(settingItem)
	if err != validator.ErrSettingEmpty {
		t.Errorf("expect err to be %v, got %v", validator.ErrSettingEmpty, err)
	}

	settingItem.Valid = true
	settingItem.Int16 = 23
	err = validator.ValidateSettingItem(settingItem)
	if err != validator.ErrSettingInvalid {
		t.Errorf("expect err to be %v, got %v", validator.ErrSettingInvalid, err)
	}
}

func TestParseSettingFromJsonAndValidate(t *testing.T) {
	jsonData := []byte(`{"salary":1,"entrepreneurship_opportunity":30}`)
	var settings models.QuestionnareSetting
	err := json.Unmarshal(jsonData, &settings)
	if err != nil {
		t.Errorf("expect no err Unmarshal json, got: %v", err)
	}
	err = validator.ValidateSettingItem(settings.Salary)
	if err != nil {
		t.Errorf("expect no err, got: %v", err)
	}

	err = validator.ValidateSettingItem(settings.TotalOpenJobs)
	if err != validator.ErrSettingEmpty {
		t.Errorf("expect err to be %v, got %v", validator.ErrSettingEmpty, err)
	}

	err = validator.ValidateSettingItem(settings.EntrepreneurshipOpportunity)
	if err != validator.ErrSettingInvalid {
		t.Errorf("expect err to be %v, got %v", validator.ErrSettingInvalid, err)
	}
}
