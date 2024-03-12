package repositories

import (
	"github.com/albugowy15/api-double-track/internal/pkg/db"
	"github.com/albugowy15/api-double-track/internal/pkg/models"
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
	_, err := db.GetDb().Exec(
		"INSERT INTO questionnare_settings (alternative_id, school_id, total_open_jobs, entrepreneurship_opportunity, salary) VALUES ($1, $2, $3, $4, $5)",
		data.AlternativeId,
		data.SchoolId,
		data.TotalOpenJobs,
		data.EntrepreneurshipOpportunity,
		data.Salary,
	)
	return err
}
