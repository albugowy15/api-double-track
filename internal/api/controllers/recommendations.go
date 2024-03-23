package controllers

import (
	"net/http"

	"github.com/albugowy15/api-double-track/internal/pkg/utils/httputil"
)

func GetRecommendations(w http.ResponseWriter, r *http.Request) {
	// by school id
	httputil.SendMessage(w, "rekomendasi", http.StatusOK)
}
