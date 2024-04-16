package ahp

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/albugowy15/api-double-track/internal/models"
)

const (
	TotalCriteria    = 5
	TotalSubCriteria = 4
	TotalAlternative = 7
)

var (
	AnswerToFloat = map[string]float32{
		"9":   9.0,
		"7":   7.0,
		"5":   5.0,
		"3":   3.0,
		"1":   1.0,
		"1/3": 1.0 / 3.0,
		"1/5": 1.0 / 5.0,
		"1/7": 1.0 / 7.0,
		"1/9": 1.0 / 9.0,
	}
	NumberToRow = map[int]int{
		15: 1,
		16: 2,
		17: 3,
		18: 4,
		19: 2,
		20: 3,
		21: 4,
		22: 3,
		23: 4,
		24: 4,
	}
	NumberToCol = map[int]int{
		15: 0,
		16: 0,
		17: 0,
		18: 0,
		19: 1,
		20: 1,
		21: 1,
		22: 2,
		23: 2,
		24: 3,
	}
	RandomIndex      = []float32{0.0, 0.0, 0.58, 0.90, 1.12, 1.24, 1.32, 1.41, 1.45, 1.49}
	AlternativeToRow = map[string]int{
		"Multimedia":                      0,
		"Teknik Elektro":                  1,
		"Teknik Listrik":                  2,
		"Tata Busana":                     3,
		"Tata Boga":                       4,
		"Tata Kecantikan":                 5,
		"Teknik Kendararaan Ringan/Motor": 6,
	}
	CriteriaToCol = map[string]int{
		"total_open_jobs":              0,
		"salary":                       1,
		"entrepreneurship_opportunity": 2,
		"interest":                     3,
		"supporting_facilites":         4,
	}
)

type (
	SubMPC            = [TotalSubCriteria][TotalSubCriteria]float32
	MPC               = [TotalCriteria][TotalCriteria]float32
	ColSum            = [TotalCriteria]float32
	SubColSum         = [TotalSubCriteria]float32
	CriteriaWeight    = [TotalCriteria]float32
	SubCriteriaWeight = [TotalSubCriteria]float32
	WeightedSum       = [TotalCriteria]float32
	DecisionMatrix    = [TotalAlternative][TotalCriteria]float32
)

func SumRow(row [TotalCriteria]float32) float32 {
	var sum float32 = 0.0
	for _, val := range row {
		sum += val
	}
	return sum
}

func SubSumRow(row [TotalSubCriteria]float32) float32 {
	var sum float32 = 0.0
	for _, val := range row {
		sum += val
	}
	return sum
}

func flipAnswer(answer string) float32 {
	if strings.ContainsAny(answer, "/") {
		split := strings.Split(answer, "/")
		return AnswerToFloat[split[1]]
	}
	return 1.0 / AnswerToFloat[answer]
}

func BuildSubMPC() SubMPC {
	subMpc := SubMPC{}
	for row := range TotalSubCriteria {
		for col := range TotalSubCriteria {
			subMpc[row][col] = (float32(TotalSubCriteria) - float32(col)) / (float32(TotalSubCriteria) - float32(row))
		}
	}
	return subMpc
}

func BuildCriteriaMPC(data []models.SubmitAnswerRequest) MPC {
	mpc := MPC{}
	// set all diagonal cell to 1.0
	for i := 0; i < TotalCriteria; i++ {
		mpc[i][i] = 1.0
	}

	comparisonQuestionStart := 15
	comparisonQuestionEnd := 24
	for _, item := range data {
		isComparisonQuestion := item.Number >= comparisonQuestionStart && item.Number <= comparisonQuestionEnd
		if isComparisonQuestion {
			row := NumberToRow[item.Number]
			col := NumberToCol[item.Number]
			mpc[row][col] = AnswerToFloat[item.Answer]
			mpc[col][row] = flipAnswer(item.Answer)
		}
	}
	return mpc
}

func CalculateColSum(mpc MPC) ColSum {
	colSum := ColSum{}
	for col := range TotalCriteria {
		for row := range TotalCriteria {
			colSum[col] += mpc[row][col]
		}
	}
	return colSum
}

func CalculateSubColSum(subMpc SubMPC) SubColSum {
	subColSum := SubColSum{}
	for col := range TotalSubCriteria {
		for row := range TotalSubCriteria {
			subColSum[col] += subMpc[row][col]
		}
	}
	return subColSum
}

func NormalizeMPC(mpc MPC, colSum ColSum) MPC {
	normalizedMpc := MPC{}
	for row := range TotalCriteria {
		for col := range TotalCriteria {
			normalizedMpc[row][col] = mpc[row][col] / colSum[col]
		}
	}
	return normalizedMpc
}

func NormalizeSubMPC(subMpc SubMPC, subColSum SubColSum) SubMPC {
	subNormMpc := SubMPC{}
	for row := range TotalSubCriteria {
		for col := range TotalSubCriteria {
			subNormMpc[row][col] = subMpc[row][col] / subColSum[col]
		}
	}
	return subNormMpc
}

func CalculateCriteriaWeight(normMpc MPC) CriteriaWeight {
	criteriaWeight := CriteriaWeight{}
	for row := range TotalCriteria {
		criteriaWeight[row] = SumRow(normMpc[row]) / TotalCriteria
	}
	return criteriaWeight
}

func CalculateSubCriteriaWeight(subNormMpc SubMPC) SubCriteriaWeight {
	subCriteriaWeight := SubCriteriaWeight{}
	for row := range TotalSubCriteria {
		subCriteriaWeight[row] = SubSumRow(subNormMpc[row]) / TotalSubCriteria
	}
	return subCriteriaWeight
}

func CalculateWeightedSum(mpc MPC, criteriaWeight WeightedSum) WeightedSum {
	weightedMpc := MPC{}
	log.Println("weightedMpc")
	for row := range TotalCriteria {
		for col := range TotalCriteria {
			weightedMpc[row][col] = mpc[row][col] * criteriaWeight[col]
		}
		log.Println(weightedMpc[row])
	}
	weightedSum := CriteriaWeight{}
	for row := range TotalCriteria {
		weightedSum[row] = SumRow(weightedMpc[row])
	}
	return weightedSum
}

func CalculateLambdaMax(weightedSum WeightedSum, criteriaWeight CriteriaWeight) float32 {
	ratios := [TotalCriteria]float32{}
	log.Println("ratios")
	for row := range TotalCriteria {
		ratios[row] = weightedSum[row] / criteriaWeight[row]
		log.Println("ratio:", ratios[row])
	}
	sum := SumRow(ratios)
	lambdaMax := sum / TotalCriteria
	return lambdaMax
}

func ConsistencyIndex(lambdaMax float32) float32 {
	return (lambdaMax - TotalCriteria) / (TotalCriteria - 1)
}

func ConsistencyRatio(consistencyIndex float32) float32 {
	return consistencyIndex / RandomIndex[TotalCriteria-1]
}

func IsAnswerConsistent(consistencyRatio float32) bool {
	return consistencyRatio < 0.10
}

func BuildDecisionMatrix(
	settings []models.QuestionnareSettingAlternative,
	answers []models.SubmitAnswerRequest,
) (DecisionMatrix, error) {
	decisionMatrix := DecisionMatrix{}
	// build from settings
	for _, setting := range settings {
		row, ok := AlternativeToRow[setting.Alternative]
		if !ok {
			message := fmt.Sprintf("no matching row for alternative: %s", setting.Alternative)
			return decisionMatrix, errors.New(message)
		}
		if !setting.TotalOpenJobs.Valid {
			return decisionMatrix, errors.New("total open jobs setting is null")
		}
		decisionMatrix[row][CriteriaToCol["total_open_jobs"]] = float32(setting.TotalOpenJobs.Int16)

		if !setting.Salary.Valid {
			return decisionMatrix, errors.New("salary setting is null")
		}
		decisionMatrix[row][CriteriaToCol["salary"]] = float32(setting.Salary.Int16)

		if !setting.EntrepreneurshipOpportunity.Valid {
			return decisionMatrix, errors.New("entrepreneurship opportunity setting is null")
		}
		decisionMatrix[row][CriteriaToCol["entrepreneurship_opportunity"]] = float32(setting.EntrepreneurshipOpportunity.Int16)
	}

	interestQuestionNumberStart := 1
	interestQuestionNumberEnd := 7

	supportingFacilitiesQuestionNumberStart := 8
	supportingFacilitiesQuestionNumberEnd := 14

	// build from answers
	for _, answer := range answers {
		isInterestQuestion := answer.Number >= interestQuestionNumberStart && answer.Number <= interestQuestionNumberEnd
		isSupportingFacilitiesQuestion := answer.Number >= supportingFacilitiesQuestionNumberStart && answer.Number <= supportingFacilitiesQuestionNumberEnd

		if isInterestQuestion {
			col := CriteriaToCol["interest"]
			row := answer.Number - interestQuestionNumberStart
			answerInt, err := strconv.Atoi(answer.Answer)
			if err != nil {
				return decisionMatrix, err
			}
			decisionMatrix[row][col] = float32(answerInt)
		} else if isSupportingFacilitiesQuestion {
			col := CriteriaToCol["supporting_facilites"]
			row := answer.Number - supportingFacilitiesQuestionNumberStart
			answerInt, err := strconv.Atoi(answer.Answer)
			if err != nil {
				return decisionMatrix, err
			}
			decisionMatrix[row][col] = float32(answerInt)
		}
	}
	return decisionMatrix, nil
}
