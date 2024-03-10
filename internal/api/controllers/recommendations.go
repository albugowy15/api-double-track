package controllers

import (
	"net/http"

	"github.con/albugowy15/api-double-track/internal/pkg/models"
	"github.con/albugowy15/api-double-track/internal/pkg/utils"
)

func GetRecommendations(w http.ResponseWriter, r *http.Request) {
	// by school id
	res := models.MessageResponse{
		Message: "rekomendasi",
	}
	utils.SendJson(w, res, http.StatusOK)
}
