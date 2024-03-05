package repositories

import (
	"github.con/albugowy15/api-double-track/internal/pkg/db"
	"github.con/albugowy15/api-double-track/internal/pkg/models"
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
