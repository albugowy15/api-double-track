package weightmethods

import (
	"net/http"

	"github.com/albugowy15/api-double-track/internal/models"
	"github.com/albugowy15/api-double-track/pkg/ahp"
)

type AHP_Weight struct {
	Total_open_jobs            float32
	Salary                     float32
	Entrepreneur_Opportunities float32
	Interest                   float32
	Supporting_Facilities      float32
}

func AhpWeight(r *http.Request, body []models.SubmitAnswerRequest, criteriaWeight ahp.CriteriaWeight) (error, AHP_Weight) {
	return nil, AHP_Weight{
		Total_open_jobs:            criteriaWeight[0],
		Salary:                     criteriaWeight[1],
		Entrepreneur_Opportunities: criteriaWeight[2],
		Interest:                   criteriaWeight[3],
		Supporting_Facilities:      criteriaWeight[4],
	}
}
