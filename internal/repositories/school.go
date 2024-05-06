package repositories

import (
	"log"

	"github.com/albugowy15/api-double-track/db"
	"github.com/albugowy15/api-double-track/internal/models"
)

func GetSchools() ([]models.School, error) {
	var schools []models.School
	err := db.AppDB.Select(&schools, `SELECT id, name FROM schools ORDER BY name ASC`)
	if err != nil {
		log.Println("db err:", err)
	}
	return schools, err
}

func GetSchoolById(schoolId string) (models.School, error) {
	school := models.School{}
	err := db.AppDB.Get(&school, "SELECT id, name FROM schools WHERE id = $1", schoolId)
	if err != nil {
		log.Println("db err:", err)
	}
	return school, err
}

func GetSchoolByStudentId(studentId string) (models.School, error) {
	school := models.School{}
	err := db.AppDB.Get(&school, "SELECT sc.id, name FROM schools sc INNER JOIN students st ON sc.id = st.school_id WHERE st.id = $1", studentId)
	if err != nil {
		log.Println("db err:", err)
	}
	return school, err
}

func GetSchoolByAdminId(adminId string) (models.School, error) {
	school := models.School{}
	err := db.AppDB.Get(&school, "SELECT s.id, name FROM schools s INNER JOIN admins a ON s.id = a.school_id WHERE a.id = $1", adminId)
	if err != nil {
		log.Println("db err:", err)
	}
	return school, err
}
