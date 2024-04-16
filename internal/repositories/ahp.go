package repositories

import (
	"log"

	"github.com/albugowy15/api-double-track/db"
	"github.com/albugowy15/api-double-track/internal/models"
	"github.com/jmoiron/sqlx"
)

type AHPRepository struct{}

var ahpRepository *AHPRepository

func GetAHPRepository() *AHPRepository {
	if ahpRepository == nil {
		ahpRepository = &AHPRepository{}
	}
	return ahpRepository
}

func (r *AHPRepository) GetAHPByStudentId(studentId string) (models.AHP, error) {
	ahp := models.AHP{}
	err := db.AppDB.Get(&ahp, "SELECT id, student_id, consistency_ratio FROM ahp WHERE student_id = $1", studentId)
	if err != nil {
		log.Println("db err:", err)
	}
	return ahp, err
}

func (r *AHPRepository) SaveAHP(data models.AHP) error {
	_, err := db.AppDB.Exec("INSERT INTO ahp (student_id, consistency_ratio) VALUES ($1, $2)", data.StudentId, data.ConsistencyRatio)
	if err != nil {
		log.Println("db err:", err)
		return err
	}
	return nil
}

func (r *AHPRepository) SaveAHPTx(data models.AHP, tx *sqlx.Tx) (int32, error) {
	var lastInsertedId int32
	err := tx.QueryRowx("INSERT INTO ahp (student_id, consistency_ratio) VALUES ($1, $2) RETURNING id", data.StudentId, data.ConsistencyRatio).Scan(&lastInsertedId)
	if err != nil {
		log.Println("db err:", err)
		return lastInsertedId, err
	}

	return lastInsertedId, nil
}

func (r *AHPRepository) SaveAHPAlternative(data models.AHPToAlternatives) error {
	_, err := db.AppDB.Exec("INSERT INTO ahp_to_alternatives (score, ahp_id, alternative_id) VALUES ($1, $2, $3)", data.Score, data.AhpId, data.AlternativeId)
	if err != nil {
		log.Println("db err:", err)
		return err
	}
	return nil
}

func (r *AHPRepository) SaveAHPAlternativeTx(data models.AHPToAlternatives, tx *sqlx.Tx) error {
	_, err := tx.Exec("INSERT INTO ahp_to_alternatives (score, ahp_id, alternative_id) VALUES ($1, $2, $3)", data.Score, data.AhpId, data.AlternativeId)
	if err != nil {
		log.Println("db err:", err)
		return err
	}
	return nil
}

func (r *AHPRepository) SaveAHPAlternatives(data []models.AHPToAlternatives) error {
	_, err := db.AppDB.NamedExec("INSERT INTO ahp_to_alternatives (score, ahp_id, alternative_id) VALUES (:score, :ahp_id, :alternative_id)", data)
	if err != nil {
		log.Println("db err:", err)
		return err
	}
	return nil
}

func (r *AHPRepository) SaveAHPAlternativesTx(data []models.AHPToAlternatives, tx *sqlx.Tx) error {
	_, err := tx.NamedExec("INSERT INTO ahp_to_alternatives (score, ahp_id, alternative_id) VALUES (:score, :ahp_id, :alternative_id)", data)
	if err != nil {
		log.Println("db err:", err)
		return err
	}
	return nil
}
