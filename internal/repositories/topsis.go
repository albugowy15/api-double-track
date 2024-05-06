package repositories

import (
	"log"

	"github.com/albugowy15/api-double-track/internal/models"
	"github.com/jmoiron/sqlx"
)

func SaveTOPSISAlternativeTx(data models.TOPSISToAlternatives, tx *sqlx.Tx) error {
	_, err := tx.Exec("INSERT INTO topsis_to_alternatives (score, topsis_id, alternative_id) VALUES ($1, $2, $3)", data.Score, data.TopsisId, data.AlternativeId)
	if err != nil {
		log.Println("db err:", err)
		return err
	}
	return nil
}

func SaveTOPSISAHPAlternativeTx(data models.TOPSISToAlternatives, tx *sqlx.Tx) error {
	_, err := tx.Exec("INSERT INTO topsis_ahp_to_alternatives (score, topsis_id, alternative_id) VALUES ($1, $2, $3)", data.Score, data.TopsisId, data.AlternativeId)
	if err != nil {
		log.Println("db err:", err)
		return err
	}
	return nil
}

func SaveTOPSISTx(data models.TOPSIS, tx *sqlx.Tx) (int32, error) {
	var lastInsertedId int32
	err := tx.QueryRowx(`
	INSERT INTO topsis (student_id)	VALUES ($1) RETURNING id`,
		data.StudentId,
	).Scan(&lastInsertedId)
	if err != nil {
		log.Println("db err : ", err)
		return lastInsertedId, err
	}
	return lastInsertedId, nil
}
