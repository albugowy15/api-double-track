package repositories

import (
	"github.com/albugowy15/api-double-track/db"
	"github.com/albugowy15/api-double-track/internal/models"
)

type SchoolRepository struct{}

var schoolRepository *SchoolRepository

func GetSchoolRepository() *SchoolRepository {
	if schoolRepository == nil {
		schoolRepository = &SchoolRepository{}
	}
	return schoolRepository
}

func (r *SchoolRepository) GetSchoolById(schoolId string) (models.School, error) {
	school := models.School{}
	err := db.AppDB.Get(&school, "SELECT id, name FROM schools WHERE id = $1", schoolId)
	return school, err
}

func (r *SchoolRepository) GetSchoolByStudentId(studentId string) (models.School, error) {
	school := models.School{}
	err := db.AppDB.Get(&school, "SELECT sc.id, name FROM schools sc INNER JOIN students st ON sc.id = st.school_id WHERE st.id = $1", studentId)
	return school, err
}

func (r *SchoolRepository) GetSchoolByAdminId(adminId string) (models.School, error) {
	school := models.School{}
	err := db.AppDB.Get(&school, "SELECT s.id, name FROM schools s INNER JOIN admins a ON s.id = a.school_id WHERE a.id = $1", adminId)
	return school, err
}
