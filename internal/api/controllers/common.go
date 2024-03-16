package controllers

import (
	"net/http"

	"github.com/albugowy15/api-double-track/internal/pkg/models"
	"github.com/albugowy15/api-double-track/internal/pkg/repositories"
	"github.com/albugowy15/api-double-track/internal/pkg/utils"
)

func GetStatistics(w http.ResponseWriter, r *http.Request) {
	res := models.Statistic{
		RegisteredStudents:       123,
		QuestionnareCompleted:    450,
		RecommendationAcceptance: 90.34,
		ConsistencyAvg:           92.54,
	}
	utils.SendJson(w, res, http.StatusOK)
}

func GetAlternatives(w http.ResponseWriter, r *http.Request) {
	s := repositories.GetAlternativeRepository()
	alternatives, err := s.GetAlternatives()
	if err != nil {
		utils.SendError(w, "internal server error", http.StatusInternalServerError)
		return
	}
	utils.SendJson(w, alternatives, http.StatusOK)
}

func GetSchool(w http.ResponseWriter, r *http.Request) {
	utils.SendError(w, "heee", http.StatusInternalServerError)
}
