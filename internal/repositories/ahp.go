package repositories

import (
	"database/sql"
	"log"

	"github.com/albugowy15/api-double-track/db"
	"github.com/albugowy15/api-double-track/internal/models"
	"github.com/jmoiron/sqlx"
)

func GetAHPByStudentId(studentId string) (models.AHP, error) {
	ahp := models.AHP{}
	err := db.AppDB.Get(&ahp, "SELECT id, student_id, consistency_ratio FROM ahp WHERE student_id = $1", studentId)
	if err != nil {
		log.Println("db err:", err)
	}
	return ahp, err
}

func SaveAHP(data models.AHP) error {
	_, err := db.AppDB.Exec("INSERT INTO ahp (student_id, consistency_ratio) VALUES ($1, $2)", data.StudentId, data.ConsistencyRatio)
	if err != nil {
		log.Println("db err:", err)
		return err
	}
	return nil
}

func SaveAHPTx(data models.AHP, tx *sqlx.Tx) (int32, error) {
	var lastInsertedId int32
	err := tx.QueryRowx("INSERT INTO ahp (student_id, consistency_ratio) VALUES ($1, $2) RETURNING id", data.StudentId, data.ConsistencyRatio).Scan(&lastInsertedId)
	if err != nil {
		log.Println("db err:", err)
		return lastInsertedId, err
	}

	return lastInsertedId, nil
}

func SaveAHPAlternative(data models.AHPToAlternatives) error {
	_, err := db.AppDB.Exec("INSERT INTO ahp_to_alternatives (score, ahp_id, alternative_id) VALUES ($1, $2, $3)", data.Score, data.AhpId, data.AlternativeId)
	if err != nil {
		log.Println("db err:", err)
		return err
	}
	return nil
}

func SaveAHPAlternativeTx(data models.AHPToAlternatives, tx *sqlx.Tx) error {
	_, err := tx.Exec("INSERT INTO ahp_to_alternatives (score, ahp_id, alternative_id) VALUES ($1, $2, $3)", data.Score, data.AhpId, data.AlternativeId)
	if err != nil {
		log.Println("db err:", err)
		return err
	}
	return nil
}

func SaveAHPAlternatives(data []models.AHPToAlternatives) error {
	_, err := db.AppDB.NamedExec("INSERT INTO ahp_to_alternatives (score, ahp_id, alternative_id) VALUES (:score, :ahp_id, :alternative_id)", data)
	if err != nil {
		log.Println("db err:", err)
		return err
	}
	return nil
}

func SaveAHPAlternativesTx(data []models.AHPToAlternatives, tx *sqlx.Tx) error {
	_, err := tx.NamedExec("INSERT INTO ahp_to_alternatives (score, ahp_id, alternative_id) VALUES (:score, :ahp_id, :alternative_id)", data)
	if err != nil {
		log.Println("db err:", err)
		return err
	}
	return nil
}

func GetAvgConsistencyRatio(schoolId string) (sql.NullFloat64, error) {
	var avgConsistencyRatio sql.NullFloat64
	err := db.AppDB.Get(
		&avgConsistencyRatio,
		`SELECT AVG(ahp.consistency_ratio) as avg_consistency_ratio FROM ahp 
    INNER JOIN students st ON st.id = ahp.student_id 
    WHERE st.school_id = $1`,
		schoolId,
	)
	if err != nil {
		log.Println("db err:", err)
	}
	return avgConsistencyRatio, err
}

func GetTotalCompleteQuestionnare(schoolId string) (int64, error) {
	var totalCompleteQestionnare int64
	err := db.AppDB.Get(
		&totalCompleteQestionnare,
		`SELECT COUNT(*) as total_complete_questionnare FROM students s 
    INNER JOIN ahp ON s.id = ahp.student_id 
    WHERE s.school_id = $1`, schoolId)
	if err != nil {
		log.Println(err)
	}
	return totalCompleteQestionnare, err
}
