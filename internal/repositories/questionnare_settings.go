package repositories

import (
	"log"
	"math"

	"github.com/albugowy15/api-double-track/db"
	"github.com/albugowy15/api-double-track/internal/models"
)

func AddQuestionnareSetting(data models.QuestionnareSetting) error {
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

func GetMissingSettings(schoolId string) ([]models.Alternative, error) {
	alternatives := []models.Alternative{}
	err := db.AppDB.Select(
		&alternatives,
		"SELECT a.id, a.alternative, a.description FROM alternatives a WHERE a.id NOT IN (SELECT alternative_id FROM questionnare_settings WHERE school_id = $1)",
		schoolId)
	return alternatives, err
}

func GetQuestionnareSettings(schoolId string) ([]models.QuestionnareSettingAlternative, error) {
	settings := []models.QuestionnareSettingAlternative{}
	err := db.AppDB.Select(
		&settings,
		"SELECT qs.id, a.alternative, qs.total_open_jobs, qs.salary, qs.entrepreneurship_opportunity FROM questionnare_settings qs INNER JOIN alternatives a ON qs.alternative_id = a.id WHERE qs.school_id = $1",
		schoolId,
	)
	return settings, err
}

func GetSumTotalOpenJobs(schoolId string) (int64, error) {
	var sum int64
	err := db.AppDB.Get(
		&sum,
		`SELECT SUM(total_open_jobs) FROM questionnare_settings
		WHERE school_id = $1`,
		schoolId,
	)

	return sum, err
}

func GetSumPowerOfTotalOpenJobs(schoolId string) (float32, error) {
	var sqrt_sum float32
	err := db.AppDB.Get(
		&sqrt_sum,
		`SELECT SUM(total_open_jobs * total_open_jobs) FROM questionnare_settings
		WHERE school_id = $1`,
		schoolId,
	)
	if err != nil {
		return 0, err
	}
	result := math.Sqrt(float64(sqrt_sum))
	return float32(result), nil
}

func GetTotalOpenJobsPerSchool(schoolID string) ([]float32, error) {
	var totalOpenJobs []float32
	err := db.AppDB.Select(
		&totalOpenJobs,
		`SELECT total_open_jobs FROM questionnare_settings
		WHERE school_id = $1`,
		schoolID,
	)

	return totalOpenJobs, err
}

func GetSumEntrepreneurshipOpportunities(schoolId string) (int64, error) {
	var sum int64
	err := db.AppDB.Get(
		&sum,
		`SELECT SUM(entrepreneurship_opportunity) FROM questionnare_settings
		WHERE school_id = $1`,
		schoolId,
	)

	return sum, err
}

func GetSumPowerOfEntrepreneurshipOpportunities(schoolId string) (float32, error) {
	var sqrt_sum float32
	err := db.AppDB.Get(
		&sqrt_sum,
		`SELECT SUM(entrepreneurship_opportunity * entrepreneurship_opportunity) FROM questionnare_settings
		WHERE school_id = $1`,
		schoolId,
	)
	if err != nil {
		return 0, err
	}
	result := math.Sqrt(float64(sqrt_sum))
	return float32(result), nil
}

func GetEntrepreneurshipOpportunitiesPerSchool(schoolID string) ([]float32, error) {
	var EntrepreneurshipOpportunities []float32
	err := db.AppDB.Select(
		&EntrepreneurshipOpportunities,
		`SELECT entrepreneurship_opportunity FROM questionnare_settings
		WHERE school_id = $1`,
		schoolID,
	)

	return EntrepreneurshipOpportunities, err
}

func GetSumSalary(schoolId string) (int64, error) {
	var sum int64
	err := db.AppDB.Get(
		&sum,
		`SELECT SUM(salary) FROM questionnare_settings
		WHERE school_id = $1`,
		schoolId,
	)

	return sum, err
}

func GetSumPowerOfSalaries(schoolId string) (float32, error) {
	var sqrt_sum float32
	err := db.AppDB.Get(
		&sqrt_sum,
		`SELECT SUM(salary * salary) FROM questionnare_settings
		WHERE school_id = $1`,
		schoolId,
	)
	if err != nil {
		return 0, err
	}
	result := math.Sqrt(float64(sqrt_sum))
	return float32(result), nil
}

func GetSalariesPerSchool(schoolID string) ([]float32, error) {
	var Salaries []float32
	err := db.AppDB.Select(
		&Salaries,
		`SELECT salary FROM questionnare_settings
		WHERE school_id = $1`,
		schoolID,
	)

	return Salaries, err
}
