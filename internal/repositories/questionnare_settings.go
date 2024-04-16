package repositories

import (
	"log"

	"github.com/albugowy15/api-double-track/db"
	"github.com/albugowy15/api-double-track/internal/models"
)

type QuestionnareSettingRepository struct{}

var questionnareSettingRepository *QuestionnareSettingRepository

func GetQuestionnareSettingRepository() *QuestionnareSettingRepository {
	if questionnareSettingRepository == nil {
		questionnareSettingRepository = &QuestionnareSettingRepository{}
	}
	return questionnareSettingRepository
}

func (r *QuestionnareSettingRepository) AddQuestionnareSetting(data models.QuestionnareSetting) error {
	tx, err := db.AppDB.Beginx()
	if err != nil {
		log.Fatalf("err start transaction: %v", err)
		return err
	}
	_, err = tx.Exec("DELETE FROM questionnare_settings WHERE alternative_id = $1 AND school_id = $2", data.AlternativeId, data.SchoolId)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(
		"INSERT INTO questionnare_settings (alternative_id, school_id, total_open_jobs, entrepreneurship_opportunity, salary) VALUES ($1, $2, $3, $4, $5)",
		data.AlternativeId,
		data.SchoolId,
		data.TotalOpenJobs,
		data.EntrepreneurshipOpportunity,
		data.Salary,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (r *QuestionnareSettingRepository) GetMissingSettings(schoolId string) ([]models.Alternative, error) {
	alternatives := []models.Alternative{}
	err := db.AppDB.Select(
		&alternatives,
		"SELECT a.id, a.alternative, a.description FROM alternatives a WHERE a.id NOT IN (SELECT alternative_id FROM questionnare_settings WHERE school_id = $1)",
		schoolId)
	return alternatives, err
}

func (r *QuestionnareSettingRepository) GetQuestionnareSettings(schoolId string) ([]models.QuestionnareSettingAlternative, error) {
	settings := []models.QuestionnareSettingAlternative{}
	err := db.AppDB.Select(
		&settings,
		"SELECT qs.id, a.alternative, qs.total_open_jobs, qs.salary, qs.entrepreneurship_opportunity FROM questionnare_settings qs INNER JOIN alternatives a ON qs.alternative_id = a.id WHERE qs.school_id = $1",
		schoolId,
	)
	return settings, err
}
