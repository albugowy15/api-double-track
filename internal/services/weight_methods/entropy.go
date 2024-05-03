package weightmethods

import (
	"math"
	"net/http"
	"strconv"

	"github.com/albugowy15/api-double-track/internal/models"
	"github.com/albugowy15/api-double-track/internal/repositories"
	"github.com/albugowy15/api-double-track/pkg/auth"
	"github.com/albugowy15/api-double-track/pkg/topsis"
)

var (
	Weight_interest                       float64
	Weight_facilities                     float64
	Weight_total_open_jobs                float64
	Weight_salaries                       float64
	Weight_entrepreneurship_opportunities float64
)

func CalculateEntropy(r *http.Request, body []models.SubmitAnswerRequest) error {
	var arr_sum [topsis.TotalCriteria]float32
	for _, answer := range body {
		is_interest_question := answer.Number >= 1 && answer.Number <= 7
		is_facilities_question := answer.Number >= 8 && answer.Number <= 14
		if is_interest_question {
			float_answer, _ := strconv.ParseFloat(answer.Answer, 32)
			arr_sum[0] += float32(float_answer)
		}
		if is_facilities_question {
			float_answer, _ := strconv.ParseFloat(answer.Answer, 32)
			arr_sum[1] += float32(float_answer)
		}
	}
	school_id_claim, _ := auth.GetJwtClaim(r, "school_id")
	school_id := school_id_claim.(string)
	sum_total_open_jobs, err := repositories.GetSumTotalOpenJobs(school_id)
	if err != nil {
		return err
	}
	arr_sum[2] = float32(sum_total_open_jobs)

	sum_entrepreneurship_opportunity, err := repositories.GetSumEntrepreneurshipOpportunities(school_id)
	if err != nil {
		return err
	}
	arr_sum[3] = float32(sum_entrepreneurship_opportunity)

	sum_salary, err := repositories.GetSumSalary(school_id)
	if err != nil {
		return err
	}
	arr_sum[4] = float32(sum_salary)

	// normalisasi
	var interest_norm [topsis.TotalAlternative]float32
	var facilities_norm [topsis.TotalAlternative]float32
	for _, answer := range body {
		if answer.Number >= 1 && answer.Number <= 7 {
			float_answer, _ := strconv.ParseFloat(answer.Answer, 32)
			interest_norm[answer.Number-1] = float32(float_answer) / arr_sum[0]
		}
		if answer.Number >= 8 && answer.Number <= 14 {
			float_answer, _ := strconv.ParseFloat(answer.Answer, 32)
			facilities_norm[answer.Number-8] = float32(float_answer) / arr_sum[1]
		}
	}

	// for i, answer := range interest_norm {
	// 	fmt.Printf("normalisasi ke %d: %v\n", i, answer)
	// }

	var TotalOpenJobs_norm [topsis.TotalAlternative]float32
	var Salaries_norm [topsis.TotalAlternative]float32
	var EntrepreneurshipOpportunity_norm [topsis.TotalAlternative]float32
	totalOpenJobs, err := repositories.GetTotalOpenJobsPerSchool(school_id)
	if err != nil {
		return err
	}
	for i, job := range totalOpenJobs {
		if i >= len(TotalOpenJobs_norm) {
			break
		}
		TotalOpenJobs_norm[i] = job
	}

	TotalOpenJobs_norm[0] = TotalOpenJobs_norm[0] / arr_sum[2]
	TotalOpenJobs_norm[1] = TotalOpenJobs_norm[1] / arr_sum[2]
	TotalOpenJobs_norm[2] = TotalOpenJobs_norm[2] / arr_sum[2]
	TotalOpenJobs_norm[3] = TotalOpenJobs_norm[3] / arr_sum[2]
	TotalOpenJobs_norm[4] = TotalOpenJobs_norm[4] / arr_sum[2]
	TotalOpenJobs_norm[5] = TotalOpenJobs_norm[5] / arr_sum[2]
	TotalOpenJobs_norm[6] = TotalOpenJobs_norm[6] / arr_sum[2]

	entrepreneurshipOpportunities, err := repositories.GetEntrepreneurshipOpportunitiesPerSchool(school_id)
	if err != nil {
		return err
	}
	for i, job := range entrepreneurshipOpportunities {
		if i >= len(EntrepreneurshipOpportunity_norm) {
			break
		}
		EntrepreneurshipOpportunity_norm[i] = job
	}

	EntrepreneurshipOpportunity_norm[0] = EntrepreneurshipOpportunity_norm[0] / arr_sum[3]
	EntrepreneurshipOpportunity_norm[1] = EntrepreneurshipOpportunity_norm[1] / arr_sum[3]
	EntrepreneurshipOpportunity_norm[2] = EntrepreneurshipOpportunity_norm[2] / arr_sum[3]
	EntrepreneurshipOpportunity_norm[3] = EntrepreneurshipOpportunity_norm[3] / arr_sum[3]
	EntrepreneurshipOpportunity_norm[4] = EntrepreneurshipOpportunity_norm[4] / arr_sum[3]
	EntrepreneurshipOpportunity_norm[5] = EntrepreneurshipOpportunity_norm[5] / arr_sum[3]
	EntrepreneurshipOpportunity_norm[6] = EntrepreneurshipOpportunity_norm[6] / arr_sum[3]

	salaries, err := repositories.GetSalariesPerSchool(school_id)
	if err != nil {
		return err
	}
	for i, job := range salaries {
		if i >= len(Salaries_norm) {
			break
		}
		Salaries_norm[i] = job
	}
	Salaries_norm[0] = Salaries_norm[0] / arr_sum[4]
	Salaries_norm[1] = Salaries_norm[1] / arr_sum[4]
	Salaries_norm[2] = Salaries_norm[2] / arr_sum[4]
	Salaries_norm[3] = Salaries_norm[3] / arr_sum[4]
	Salaries_norm[4] = Salaries_norm[4] / arr_sum[4]
	Salaries_norm[5] = Salaries_norm[5] / arr_sum[4]
	Salaries_norm[6] = Salaries_norm[6] / arr_sum[4]

	/**
	ENTROPY
	*/

	var SumEntropyInterest float32
	for _, val := range interest_norm {
		SumEntropyInterest += val * float32(math.Log(float64(val)))
	}
	entropy_interest := -1.0 / math.Log(topsis.TotalAlternative) * float64(SumEntropyInterest)
	// fmt.Println("entropy interest :", entropy_interest)

	var SumEntropyFacilities float32
	for _, val := range facilities_norm {
		SumEntropyFacilities += val * float32(math.Log(float64(val)))
	}
	entropy_facilities := -1.0 / math.Log(topsis.TotalAlternative) * float64(SumEntropyFacilities)
	// fmt.Println("entropy facilities :", entropy_facilities)

	var SumEntropyTotalOpenJobs float32
	for _, val := range TotalOpenJobs_norm {
		SumEntropyTotalOpenJobs += val * float32(math.Log(float64(val)))
	}
	entropy_total_open_jobs := -1.0 / math.Log(topsis.TotalAlternative) * float64(SumEntropyTotalOpenJobs)
	// fmt.Println("entropy total open jobs :", entropy_total_open_jobs)

	var SumSalaries float32
	for _, val := range Salaries_norm {
		SumSalaries += val * float32(math.Log(float64(val)))
	}
	entropy_salaries := -1.0 / math.Log(topsis.TotalAlternative) * float64(SumSalaries)
	// fmt.Println("entropy salaries :", entropy_salaries)

	var SumEntrepreneurshipOpportunity float32
	for _, val := range EntrepreneurshipOpportunity_norm {
		SumEntrepreneurshipOpportunity += val * float32(math.Log(float64(val)))
	}
	entropy_entrepreneurship_opportunities := -1.0 / math.Log(topsis.TotalAlternative) * float64(SumEntrepreneurshipOpportunity)
	// fmt.Println("entropy entrepreneurship opportunities :", entropy_entrepreneurship_opportunities)

	/**
	DEGREE OF DIFFERENTIATION
	*/

	dod_interest := 1 - entropy_interest
	dod_facilities := 1 - entropy_facilities
	dod_total_open_jobs := 1 - entropy_total_open_jobs
	dod_salaries := 1 - entropy_salaries
	dod_entrepreneurship_opportunities := 1 - entropy_entrepreneurship_opportunities

	total_dod := (dod_interest + dod_facilities + dod_total_open_jobs + dod_salaries + dod_entrepreneurship_opportunities)
	// fmt.Println("total dod : ", total_dod)
	/**
	weights of entropy
	*/

	Weight_interest = dod_interest / total_dod
	Weight_facilities = dod_facilities / total_dod
	Weight_total_open_jobs = dod_total_open_jobs / total_dod
	Weight_salaries = dod_salaries / total_dod
	Weight_entrepreneurship_opportunities = dod_entrepreneurship_opportunities / total_dod

	// fmt.Println("weight_interest : ", Weight_interest)
	// fmt.Println("weight_facilities : ", Weight_facilities)
	// fmt.Println("weight_total_open_jobs : ", Weight_total_open_jobs)
	// fmt.Println("weight_salaries : ", Weight_salaries)
	// fmt.Println("weight_entrepreneurship_opportunities : ", Weight_entrepreneurship_opportunities)

	// fmt.Println("Total : ", Weight_interest+Weight_facilities+Weight_total_open_jobs+Weight_salaries+Weight_entrepreneurship_opportunities)
	return nil
}
