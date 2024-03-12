package repositories

import (
	"github.com/albugowy15/api-double-track/internal/pkg/db"
	"github.com/albugowy15/api-double-track/internal/pkg/models"
)

type AlternativeRepository struct{}

var alternativeRepository *AlternativeRepository

func GetAlternativeRepository() *AlternativeRepository {
	if alternativeRepository == nil {
		alternativeRepository = &AlternativeRepository{}
	}
	return alternativeRepository
}

func (r *AlternativeRepository) GetAlternatives() ([]models.Alternative, error) {
	alternatves := []models.Alternative{}
	err := db.GetDb().Select(&alternatves, "SELECT id, alternative, description FROM alternatives")
	return alternatves, err
}

func (r *AlternativeRepository) GetAlternativeByName(name string) (models.Alternative, error) {
	alternative := models.Alternative{}
	err := db.GetDb().Get(&alternative, "SELECT id, alternative, description FROM alternatives WHERE alternative = $1", name)
	return alternative, err
}

func (r *AlternativeRepository) GetAlternativeById(id int) (models.Alternative, error) {
	alternative := models.Alternative{}
	err := db.GetDb().Get(&alternative, "SELECT id, alternative, description FROM alternatives WHERE id = $1", id)
	return alternative, err
}
