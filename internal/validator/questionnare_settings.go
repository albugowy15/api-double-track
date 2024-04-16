package validator

import (
	"errors"

	"github.com/albugowy15/api-double-track/internal/models"
	"github.com/albugowy15/api-double-track/internal/repositories"
	"github.com/guregu/null/v5"
)

var (
	ErrAlternativeIdEmpty    = errors.New("alternative_id wajib diisi")
	ErrAlternativeIdNotFound = errors.New("alternative_id tidak ditemukan")
	ErrSettingInvalid        = errors.New("nilai setting invalid")
	ErrSettingEmpty          = errors.New("nilai semua setting wajib diisi")
)

func validateAlternativeId(skillId int) error {
	if skillId == 0 {
		return ErrAlternativeIdEmpty
	}
	_, err := repositories.GetAlternativeRepository().GetAlternativeById(skillId)
	if err != nil {
		return ErrAlternativeIdNotFound
	}
	return nil
}

func ValidateSettingItem(settingVal null.Int16) error {
	if !settingVal.Valid {
		return ErrSettingEmpty
	}
	if settingVal.Int16 < 1 || settingVal.Int16 > 4 {
		return ErrSettingInvalid
	}
	return nil
}

func ValidateQuestionnareSettings(data models.QuestionnareSetting) error {
	if err := validateAlternativeId(data.AlternativeId); err != nil {
		return err
	}
	if err := ValidateSettingItem(data.TotalOpenJobs); err != nil {
		return err
	}
	if err := ValidateSettingItem(data.Salary); err != nil {
		return err
	}
	if err := ValidateSettingItem(data.EntrepreneurshipOpportunity); err != nil {
		return err
	}
	return nil
}
