package repositories

import (
	"github.com/albugowy15/api-double-track/db"
	"github.com/albugowy15/api-double-track/internal/models"
)

func GetAlternatives() ([]models.Alternative, error) {
	alternatves := []models.Alternative{}
	err := db.AppDB.Select(&alternatves, "SELECT id, alternative, description FROM alternatives")
	return alternatves, err
}

func GetAlternativeByName(name string) (models.Alternative, error) {
	alternative := models.Alternative{}
	err := db.AppDB.Get(&alternative, "SELECT id, alternative, description FROM alternatives WHERE alternative = $1", name)
	return alternative, err
}

func GetAlternativeById(id int) (models.Alternative, error) {
	alternative := models.Alternative{}
	err := db.AppDB.Get(&alternative, "SELECT id, alternative, description FROM alternatives WHERE id = $1", id)
	return alternative, err
}
