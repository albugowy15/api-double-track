package weightmethods

import (
	"net/http"

	"github.com/albugowy15/api-double-track/internal/models"
	"github.com/albugowy15/api-double-track/pkg/ahp"
	"github.com/albugowy15/api-double-track/pkg/topsis"
)

type CombinativeWeight struct {
	Total_open_jobs            float32
	Salary                     float32
	Entrepreneur_Opportunities float32
	Interest                   float32
	Supporting_Facilities      float32
}

func CalculateCombinative(r *http.Request, body []models.SubmitAnswerRequest, criteriaWeight ahp.CriteriaWeight) (error, CombinativeWeight) {
	var arr_sum [topsis.TotalCriteria]float32
	_, ahp := AhpWeight(r, body, criteriaWeight)
	arr_sum[0] = ahp.Interest * float32(Weight_interest)
	arr_sum[1] = ahp.Supporting_Facilities * float32(Weight_facilities)
	arr_sum[2] = ahp.Entrepreneur_Opportunities * float32(Weight_entrepreneurship_opportunities)
	arr_sum[3] = ahp.Total_open_jobs * float32(Weight_total_open_jobs)
	arr_sum[4] = ahp.Salary * float32(Weight_salaries)
	total := arr_sum[0] + arr_sum[1] + arr_sum[2] + arr_sum[3] + arr_sum[4]
	return nil, CombinativeWeight{
		Interest:                   arr_sum[0] / total,
		Supporting_Facilities:      arr_sum[1] / total,
		Entrepreneur_Opportunities: arr_sum[2] / total,
		Total_open_jobs:            arr_sum[3] / total,
		Salary:                     arr_sum[4] / total,
	}
}
