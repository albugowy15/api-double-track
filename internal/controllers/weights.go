package controllers

import (
	"net/http"

	"github.com/albugowy15/api-double-track/internal/models"
	weightmethods "github.com/albugowy15/api-double-track/internal/services/weight_methods"
	"github.com/albugowy15/api-double-track/pkg/ahp"
	"github.com/albugowy15/api-double-track/pkg/auth"
	"github.com/albugowy15/api-double-track/pkg/httpx"
)

func HandleGetWeights(w http.ResponseWriter, r *http.Request) {
	_, _ = auth.GetJwtClaim(r, "user_id")
	entropy := models.Weights{
		Interest:                  weightmethods.Weight_interest,
		Facilities:                weightmethods.Weight_facilities,
		TotalOpenJobs:             weightmethods.Weight_total_open_jobs,
		Salaries:                  weightmethods.Weight_salaries,
		EntrepreneurOpportunities: weightmethods.Weight_entrepreneurship_opportunities,
	}
	var criteriaWeight ahp.CriteriaWeight
	ahp := models.Weights{
		TotalOpenJobs:             float64(criteriaWeight[0]),
		Salaries:                  float64(criteriaWeight[1]),
		EntrepreneurOpportunities: float64(criteriaWeight[2]),
		Interest:                  float64(criteriaWeight[3]),
		Facilities:                float64(criteriaWeight[4]),
	}
	res := models.CriteriaWeights{
		Entropy: entropy,
		AHP:     ahp,
	}
	httpx.SendData(w, res, http.StatusOK)
}
