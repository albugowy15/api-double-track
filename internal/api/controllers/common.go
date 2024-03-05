package controllers

import (
	"net/http"

	"github.con/albugowy15/api-double-track/internal/pkg/repositories"
	"github.con/albugowy15/api-double-track/internal/pkg/utils"
)

func GetStatistics(w http.ResponseWriter, r *http.Request) {
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
