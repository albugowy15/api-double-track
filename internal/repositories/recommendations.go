package repositories

import (
	"log"

	"github.com/albugowy15/api-double-track/db"
	"github.com/albugowy15/api-double-track/internal/models"
)

func GetRecommendationsBySchoolId(schoolId string) ([]models.StudentRecommendation, error) {
	students := []models.StudentRecommendation{}
	tx, err := db.AppDB.Beginx()
	if err != nil {
		log.Println("db err:", err)
		return students, err
	}

	err = tx.Select(
		&students,
		`SELECT ahp.id, ahp.consistency_ratio, s.id as student_id, s.fullname, s.nisn FROM ahp 
    INNER JOIN students s ON s.id = ahp.student_id 
    WHERE s.school_id = $1`,
		schoolId,
	)
	if err != nil {
		log.Println("db err:", err)
		tx.Rollback()
		return students, err
	}

	for idx, student := range students {
		ahpResults := []models.RecommendationResult{}
		err = tx.Select(
			&ahpResults,
			`SELECT ata.id, ata.score, a.alternative, a.description FROM ahp_to_alternatives ata 
      INNER JOIN alternatives a ON a.id = ata.alternative_id 
      WHERE ahp_id = $1
      ORDER BY ata.score DESC`,
			student.Id,
		)
		if err != nil {
			log.Println("db err:", err)
			tx.Rollback()
			return students, err
		}
		students[idx].AhpResults = ahpResults

		// do query for topsis
	}

	return students, nil
}

func GetAHPRecommendations(studentId string) ([]models.RecommendationResult, error) {
	recommendations := []models.RecommendationResult{}
	err := db.AppDB.Select(
		&recommendations,
		`SELECT ata.id, a.alternative, ata.score, a.description FROM ahp_to_alternatives ata 
    INNER JOIN alternatives a ON a.id = ata.alternative_id
    INNER JOIN ahp ON ahp.id = ata.ahp_id
    WHERE ahp.student_id = $1
    ORDER BY ata.score DESC`,
		studentId,
	)
	if err != nil {
		log.Println("db err:", err)
	}
	return recommendations, err
}

func GetAHPConsistencyRatio(studentId string) (float32, error) {
	var consistencyRatio float32
	err := db.AppDB.Get(&consistencyRatio, "SELECT consistency_ratio FROM ahp WHERE student_id = $1", studentId)
	if err != nil {
		log.Println("db err:", err)
	}
	return consistencyRatio, err
}
