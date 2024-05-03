package services

import (
	"math"
	"net/http"
	"strconv"

	"github.com/albugowy15/api-double-track/internal/models"
	"github.com/albugowy15/api-double-track/internal/repositories"
	weightmethods "github.com/albugowy15/api-double-track/internal/services/weight_methods"
	"github.com/albugowy15/api-double-track/pkg/auth"
	"github.com/albugowy15/api-double-track/pkg/httpx"
	"github.com/albugowy15/api-double-track/pkg/topsis"
	"github.com/jmoiron/sqlx"
)

/**
* perhitungan melalui switch
* case 1 : entropy (done)
* case 2 : SAW
* case 3 : Topsis
* case 4 : AHP
 */

type TOPSISServiceError struct {
	Err        error
	StatusCode int
}

func (e *TOPSISServiceError) Error() string {
	return e.Err.Error()
}

func CalculateTopsis(r *http.Request, body []models.SubmitAnswerRequest, tx *sqlx.Tx) error {
	var arr_sum [topsis.TotalCriteria]float32
	for _, answer := range body {
		is_interest_question := answer.Number >= 1 && answer.Number <= 7
		is_facilities_question := answer.Number >= 8 && answer.Number <= 14
		if is_interest_question {
			float_answer, _ := strconv.ParseFloat(answer.Answer, 32)
			squared := float32(math.Pow(float_answer, 2))
			arr_sum[0] += squared
		}
		if is_facilities_question {
			float_answer, _ := strconv.ParseFloat(answer.Answer, 32)
			squared := float32(math.Pow(float_answer, 2))
			arr_sum[1] += squared
		}
	}
	arr_sum[0] = float32(math.Sqrt(float64(arr_sum[0])))
	arr_sum[1] = float32(math.Sqrt(float64(arr_sum[1])))

	school_id_claim, _ := auth.GetJwtClaim(r, "school_id")
	school_id := school_id_claim.(string)

	sum_sqrt_total_open_jobs, err := repositories.GetSumPowerOfTotalOpenJobs(school_id)
	if err != nil {
		return err
	}
	arr_sum[2] = sum_sqrt_total_open_jobs

	sum_sqrt_entrepreneurship_opportunity, err := repositories.GetSumPowerOfEntrepreneurshipOpportunities(school_id)
	if err != nil {
		return err
	}
	arr_sum[3] = sum_sqrt_entrepreneurship_opportunity

	sum_sqrt_salaries, err := repositories.GetSumPowerOfSalaries(school_id)
	if err != nil {
		return err
	}
	arr_sum[4] = sum_sqrt_salaries

	// fmt.Println("arr_sum[0] : ", arr_sum[0]) // interest
	// fmt.Println("arr_sum[1] : ", arr_sum[1]) // facilities
	// fmt.Println("arr_sum[2] : ", arr_sum[2]) // total open jobs
	// fmt.Println("arr_sum[3] : ", arr_sum[3]) // entrepreneurship opportunity
	// fmt.Println("arr_sum[4] : ", arr_sum[4]) // salaries

	/**
	TOPSIS NORMALISASI ??
	*/
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
	weighted entropy
	*/

	for idx, val := range interest_norm {
		interest_norm[idx] = val * float32(weightmethods.Weight_interest)
		// fmt.Printf("interest_norm [%d] : %v\n", idx, interest_norm[idx])
	}

	for idx, val := range facilities_norm {
		facilities_norm[idx] = val * float32(weightmethods.Weight_facilities)
		// fmt.Printf("facilities_norm [%d] : %v\n", idx, facilities_norm[idx])
	}

	for idx, val := range TotalOpenJobs_norm {
		TotalOpenJobs_norm[idx] = val * float32(weightmethods.Weight_total_open_jobs)
		// fmt.Printf("TotalOpenJobs_norm [%d] : %v\n", idx, TotalOpenJobs_norm[idx])
	}

	for idx, val := range EntrepreneurshipOpportunity_norm {
		EntrepreneurshipOpportunity_norm[idx] = val * float32(weightmethods.Weight_entrepreneurship_opportunities)
		// fmt.Printf("EntrepreneurshipOpportunity_norm [%d] : %v\n", idx, EntrepreneurshipOpportunity_norm[idx])
	}

	for idx, val := range Salaries_norm {
		Salaries_norm[idx] = val * float32(weightmethods.Weight_salaries)
		// fmt.Printf("Salaries_norm [%d] : %v\n", idx, Salaries_norm[idx])
	}

	/*
	* PIS && NIS
	 */

	var pis_interest, nis_interest float32
	if len(interest_norm) > 0 {
		pis_interest = interest_norm[0] // Initialize pis_interest with the first element
		for i := 1; i < len(interest_norm); i++ {
			pis_interest = float32(math.Max(float64(pis_interest), float64(interest_norm[i])))
		}
		nis_interest = interest_norm[0]
		for i := 1; i < len(interest_norm); i++ {
			nis_interest = float32(math.Min(float64(nis_interest), float64(interest_norm[i])))
		}
	}

	var pis_facilities, nis_facilities float32
	if len(facilities_norm) > 0 {
		pis_facilities = facilities_norm[0]
		for i := 1; i < len(facilities_norm); i++ {
			pis_facilities = float32(math.Max(float64(pis_facilities), float64(facilities_norm[i])))
		}
		nis_facilities = facilities_norm[0]
		for i := 1; i < len(facilities_norm); i++ {
			nis_facilities = float32(math.Min(float64(nis_facilities), float64(facilities_norm[i])))
		}
	}

	var pis_total_open_jobs, nis_total_open_jobs float32
	if len(TotalOpenJobs_norm) > 0 {
		pis_total_open_jobs = TotalOpenJobs_norm[0]
		for i := 1; i < len(TotalOpenJobs_norm); i++ {
			pis_total_open_jobs = float32(math.Max(float64(pis_total_open_jobs), float64(TotalOpenJobs_norm[i])))
		}
		nis_total_open_jobs = TotalOpenJobs_norm[0]
		for i := 1; i < len(TotalOpenJobs_norm); i++ {
			nis_total_open_jobs = float32(math.Min(float64(nis_total_open_jobs), float64(TotalOpenJobs_norm[i])))
		}
	}

	var pis_entrepreneurship_opportunities, nis_entrepreneurship_opportunities float32
	if len(EntrepreneurshipOpportunity_norm) > 0 {
		pis_entrepreneurship_opportunities = EntrepreneurshipOpportunity_norm[0]
		for i := 1; i < len(EntrepreneurshipOpportunity_norm); i++ {
			pis_entrepreneurship_opportunities = float32(math.Max(float64(pis_entrepreneurship_opportunities), float64(EntrepreneurshipOpportunity_norm[i])))
		}
		nis_entrepreneurship_opportunities = EntrepreneurshipOpportunity_norm[0]
		for i := 1; i < len(EntrepreneurshipOpportunity_norm); i++ {
			nis_entrepreneurship_opportunities = float32(math.Min(float64(nis_entrepreneurship_opportunities), float64(EntrepreneurshipOpportunity_norm[i])))
		}
	}

	var pis_salary, nis_salary float32
	if len(Salaries_norm) > 0 {
		pis_salary = Salaries_norm[0]
		for i := 1; i < len(Salaries_norm); i++ {
			pis_salary = float32(math.Max(float64(pis_salary), float64(Salaries_norm[i])))
		}
		nis_salary = Salaries_norm[0]
		for i := 1; i < len(Salaries_norm); i++ {
			nis_salary = float32(math.Min(float64(nis_salary), float64(Salaries_norm[i])))
		}
	}

	// fmt.Println("pis_interest : ", pis_interest)
	// fmt.Println("nis_interest : ", nis_interest)
	// fmt.Println("pis_facilities : ", pis_facilities)
	// fmt.Println("nis_facilities : ", nis_facilities)
	// fmt.Println("pis_total_open_jobs : ", pis_total_open_jobs)
	// fmt.Println("nis_total_open_jobs : ", nis_total_open_jobs)
	// fmt.Println("pis_entrepreneurship_opportunities : ", pis_entrepreneurship_opportunities)
	// fmt.Println("nis_entrepreneurship_opportunities : ", nis_entrepreneurship_opportunities)
	// fmt.Println("pis_salary : ", pis_salary)
	// fmt.Println("nis_salary : ", nis_salary)

	/**
	Eucledian
	*/
	var d_multimedia [2]float32
	d_multimedia[0] = float32(math.Sqrt(
		(math.Pow(float64(interest_norm[0])-float64(pis_interest), 2)) +
			(math.Pow(float64(facilities_norm[0])-float64(pis_facilities), 2)) +
			(math.Pow(float64(TotalOpenJobs_norm[0])-float64(pis_total_open_jobs), 2)) +
			(math.Pow(float64(EntrepreneurshipOpportunity_norm[0])-float64(pis_entrepreneurship_opportunities), 2)) +
			(math.Pow(float64(Salaries_norm[0])-float64(pis_salary), 2))))

	d_multimedia[1] = float32(math.Sqrt(
		(math.Pow(float64(interest_norm[0])-float64(nis_interest), 2)) +
			(math.Pow(float64(facilities_norm[0])-float64(nis_facilities), 2)) +
			(math.Pow(float64(TotalOpenJobs_norm[0])-float64(nis_total_open_jobs), 2)) +
			(math.Pow(float64(EntrepreneurshipOpportunity_norm[0])-float64(nis_entrepreneurship_opportunities), 2)) +
			(math.Pow(float64(Salaries_norm[0])-float64(nis_salary), 2))))

	var d_elektro [2]float32
	d_elektro[0] = float32(math.Sqrt(
		(math.Pow(float64(interest_norm[1])-float64(pis_interest), 2)) +
			(math.Pow(float64(facilities_norm[1])-float64(pis_facilities), 2)) +
			(math.Pow(float64(TotalOpenJobs_norm[1])-float64(pis_total_open_jobs), 2)) +
			(math.Pow(float64(EntrepreneurshipOpportunity_norm[1])-float64(pis_entrepreneurship_opportunities), 2)) +
			(math.Pow(float64(Salaries_norm[1])-float64(pis_salary), 2))))

	d_elektro[1] = float32(math.Sqrt(
		(math.Pow(float64(interest_norm[1])-float64(nis_interest), 2)) +
			(math.Pow(float64(facilities_norm[1])-float64(nis_facilities), 2)) +
			(math.Pow(float64(TotalOpenJobs_norm[1])-float64(nis_total_open_jobs), 2)) +
			(math.Pow(float64(EntrepreneurshipOpportunity_norm[1])-float64(nis_entrepreneurship_opportunities), 2)) +
			(math.Pow(float64(Salaries_norm[1])-float64(nis_salary), 2))))

	var d_listrik [2]float32
	d_listrik[0] = float32(math.Sqrt(
		(math.Pow(float64(interest_norm[2])-float64(pis_interest), 2)) +
			(math.Pow(float64(facilities_norm[2])-float64(pis_facilities), 2)) +
			(math.Pow(float64(TotalOpenJobs_norm[2])-float64(pis_total_open_jobs), 2)) +
			(math.Pow(float64(EntrepreneurshipOpportunity_norm[2])-float64(pis_entrepreneurship_opportunities), 2)) +
			(math.Pow(float64(Salaries_norm[2])-float64(pis_salary), 2))))

	d_listrik[1] = float32(math.Sqrt(
		(math.Pow(float64(interest_norm[2])-float64(nis_interest), 2)) +
			(math.Pow(float64(facilities_norm[2])-float64(nis_facilities), 2)) +
			(math.Pow(float64(TotalOpenJobs_norm[2])-float64(nis_total_open_jobs), 2)) +
			(math.Pow(float64(EntrepreneurshipOpportunity_norm[2])-float64(nis_entrepreneurship_opportunities), 2)) +
			(math.Pow(float64(Salaries_norm[2])-float64(nis_salary), 2))))

	var d_busana [2]float32
	d_busana[0] = float32(math.Sqrt(
		(math.Pow(float64(interest_norm[3])-float64(pis_interest), 2)) +
			(math.Pow(float64(facilities_norm[3])-float64(pis_facilities), 2)) +
			(math.Pow(float64(TotalOpenJobs_norm[3])-float64(pis_total_open_jobs), 2)) +
			(math.Pow(float64(EntrepreneurshipOpportunity_norm[3])-float64(pis_entrepreneurship_opportunities), 2)) +
			(math.Pow(float64(Salaries_norm[3])-float64(pis_salary), 2))))

	d_busana[1] = float32(math.Sqrt(
		(math.Pow(float64(interest_norm[3])-float64(nis_interest), 2)) +
			(math.Pow(float64(facilities_norm[3])-float64(nis_facilities), 2)) +
			(math.Pow(float64(TotalOpenJobs_norm[3])-float64(nis_total_open_jobs), 2)) +
			(math.Pow(float64(EntrepreneurshipOpportunity_norm[3])-float64(nis_entrepreneurship_opportunities), 2)) +
			(math.Pow(float64(Salaries_norm[3])-float64(nis_salary), 2))))

	var d_boga [2]float32
	d_boga[0] = float32(math.Sqrt(
		(math.Pow(float64(interest_norm[4])-float64(pis_interest), 2)) +
			(math.Pow(float64(facilities_norm[4])-float64(pis_facilities), 2)) +
			(math.Pow(float64(TotalOpenJobs_norm[4])-float64(pis_total_open_jobs), 2)) +
			(math.Pow(float64(EntrepreneurshipOpportunity_norm[4])-float64(pis_entrepreneurship_opportunities), 2)) +
			(math.Pow(float64(Salaries_norm[4])-float64(pis_salary), 2))))
	d_boga[1] = float32(math.Sqrt(
		(math.Pow(float64(interest_norm[4])-float64(nis_interest), 2)) +
			(math.Pow(float64(facilities_norm[4])-float64(nis_facilities), 2)) +
			(math.Pow(float64(TotalOpenJobs_norm[4])-float64(nis_total_open_jobs), 2)) +
			(math.Pow(float64(EntrepreneurshipOpportunity_norm[4])-float64(nis_entrepreneurship_opportunities), 2)) +
			(math.Pow(float64(Salaries_norm[4])-float64(nis_salary), 2))))

	var d_kecantikan [2]float32
	d_kecantikan[0] = float32(math.Sqrt(
		(math.Pow(float64(interest_norm[5])-float64(pis_interest), 2)) +
			(math.Pow(float64(facilities_norm[5])-float64(pis_facilities), 2)) +
			(math.Pow(float64(TotalOpenJobs_norm[5])-float64(pis_total_open_jobs), 2)) +
			(math.Pow(float64(EntrepreneurshipOpportunity_norm[5])-float64(pis_entrepreneurship_opportunities), 2)) +
			(math.Pow(float64(Salaries_norm[5])-float64(pis_salary), 2))))
	d_kecantikan[1] = float32(math.Sqrt(
		(math.Pow(float64(interest_norm[5])-float64(nis_interest), 2)) +
			(math.Pow(float64(facilities_norm[5])-float64(nis_facilities), 2)) +
			(math.Pow(float64(TotalOpenJobs_norm[5])-float64(nis_total_open_jobs), 2)) +
			(math.Pow(float64(EntrepreneurshipOpportunity_norm[5])-float64(nis_entrepreneurship_opportunities), 2)) +
			(math.Pow(float64(Salaries_norm[5])-float64(nis_salary), 2))))

	var d_mesin [2]float32
	d_mesin[0] = float32(math.Sqrt(
		(math.Pow(float64(interest_norm[6])-float64(pis_interest), 2)) +
			(math.Pow(float64(facilities_norm[6])-float64(pis_facilities), 2)) +
			(math.Pow(float64(TotalOpenJobs_norm[6])-float64(pis_total_open_jobs), 2)) +
			(math.Pow(float64(EntrepreneurshipOpportunity_norm[6])-float64(pis_entrepreneurship_opportunities), 2)) +
			(math.Pow(float64(Salaries_norm[6])-float64(pis_salary), 2))))
	d_mesin[1] = float32(math.Sqrt(
		(math.Pow(float64(interest_norm[6])-float64(nis_interest), 2)) +
			(math.Pow(float64(facilities_norm[6])-float64(nis_facilities), 2)) +
			(math.Pow(float64(TotalOpenJobs_norm[6])-float64(nis_total_open_jobs), 2)) +
			(math.Pow(float64(EntrepreneurshipOpportunity_norm[6])-float64(nis_entrepreneurship_opportunities), 2)) +
			(math.Pow(float64(Salaries_norm[6])-float64(nis_salary), 2))))

	// for idx, val := range d_multimedia {
	// 	fmt.Printf("d_multimedia[%d] : %v\n", idx, val)
	// }
	// for idx, val := range d_elektro {
	// 	fmt.Printf("d_elektro[%d] : %v\n", idx, val)
	// }
	// for idx, val := range d_listrik {
	// 	fmt.Printf("d_listrik[%d] : %v\n", idx, val)
	// }
	// for idx, val := range d_busana {
	// 	fmt.Printf("d_busana[%d] : %v\n", idx, val)
	// }
	// for idx, val := range d_boga {
	// 	fmt.Printf("d_boga[%d] : %v\n", idx, val)
	// }
	// for idx, val := range d_kecantikan {
	// 	fmt.Printf("d_kecantikan[%d] : %v\n", idx, val)
	// }
	// for idx, val := range d_mesin {
	// 	fmt.Printf("d_mesin[%d] : %v\n", idx, val)
	// }
	/**
	Perankingan
	*/
	multimedia := d_multimedia[0] / (d_multimedia[0] + d_multimedia[1])
	elektro := d_elektro[0] / (d_elektro[0] + d_elektro[1])
	listrik := d_listrik[0] / (d_listrik[0] + d_listrik[1])
	busana := d_busana[0] / (d_busana[0] + d_busana[1])
	boga := d_boga[0] / (d_boga[0] + d_boga[1])
	kecantikan := d_kecantikan[0] / (d_kecantikan[0] + d_kecantikan[1])
	mesin := d_mesin[0] / (d_mesin[0] + d_mesin[1])

	// fmt.Printf("multimedia : %v\nelektro : %v\nlistrik : %v\nbusana : %v\nboga : %v\nkecantikan : %v\nmesin : %v\n",
	// 	multimedia, elektro, listrik, busana, boga, kecantikan, mesin)

	// save topsis
	studentIdClaim, _ := auth.GetJwtClaim(r, "user_id")
	studentId := studentIdClaim.(string)
	insertedId, err := repositories.SaveTOPSISTx(models.TOPSIS{
		StudentId: studentId,
	}, tx)
	if err != nil {
		return &TOPSISServiceError{
			StatusCode: http.StatusInternalServerError,
			Err:        httpx.ErrInternalServer,
		}
	}

	// 1
	multimedia_repository, err := repositories.GetAlternativeByName("Multimedia")
	if err != nil {
		return err
	}
	data := models.TOPSISToAlternatives{
		TopsisId:      insertedId,
		Score:         multimedia,
		AlternativeId: multimedia_repository.Id,
	}
	err = repositories.SaveTOPSISAlternativeTx(data, tx)
	if err != nil {
		return err
	}

	// 2
	elektro_repository, err := repositories.GetAlternativeByName("Teknik Elektro")
	if err != nil {
		return err
	}
	data = models.TOPSISToAlternatives{
		TopsisId:      insertedId,
		Score:         elektro,
		AlternativeId: elektro_repository.Id,
	}
	err = repositories.SaveTOPSISAlternativeTx(data, tx)
	if err != nil {
		return err
	}

	// 3
	listrik_repository, err := repositories.GetAlternativeByName("Teknik Listrik")
	if err != nil {
		return err
	}
	data = models.TOPSISToAlternatives{
		TopsisId:      insertedId,
		Score:         listrik,
		AlternativeId: listrik_repository.Id,
	}
	err = repositories.SaveTOPSISAlternativeTx(data, tx)
	if err != nil {
		return err
	}

	// 4
	busana_repository, err := repositories.GetAlternativeByName("Tata Busana")
	if err != nil {
		return err
	}
	data = models.TOPSISToAlternatives{
		TopsisId:      insertedId,
		Score:         busana,
		AlternativeId: busana_repository.Id,
	}
	err = repositories.SaveTOPSISAlternativeTx(data, tx)
	if err != nil {
		return err
	}

	// 5
	boga_repository, err := repositories.GetAlternativeByName("Tata Boga")
	if err != nil {
		return err
	}
	data = models.TOPSISToAlternatives{
		TopsisId:      insertedId,
		Score:         boga,
		AlternativeId: boga_repository.Id,
	}
	err = repositories.SaveTOPSISAlternativeTx(data, tx)
	if err != nil {
		return err
	}

	// 6
	kecantikan_repository, err := repositories.GetAlternativeByName("Tata Kecantikan")
	if err != nil {
		return err
	}
	data = models.TOPSISToAlternatives{
		TopsisId:      insertedId,
		Score:         kecantikan,
		AlternativeId: kecantikan_repository.Id,
	}
	err = repositories.SaveTOPSISAlternativeTx(data, tx)
	if err != nil {
		return err
	}

	// 7
	mesin_repository, err := repositories.GetAlternativeByName("Teknik Kendararaan Ringan/Motor")
	if err != nil {
		return err
	}
	data = models.TOPSISToAlternatives{
		TopsisId:      insertedId,
		Score:         mesin,
		AlternativeId: mesin_repository.Id,
	}
	err = repositories.SaveTOPSISAlternativeTx(data, tx)
	if err != nil {
		return err
	}

	return nil
}
