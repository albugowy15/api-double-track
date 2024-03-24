package services

import (
	"errors"
	"log"
	"net/http"

	"github.com/albugowy15/api-double-track/internal/pkg/ahp"
	"github.com/albugowy15/api-double-track/internal/pkg/models"
	"github.com/albugowy15/api-double-track/internal/pkg/repositories"
	"github.com/albugowy15/api-double-track/internal/pkg/utils/httputil"
	"github.com/albugowy15/api-double-track/internal/pkg/utils/jwt"
)

func CalculateAHP(r *http.Request, body []models.SubmitAnswerRequest) error {
	mpc := ahp.BuildCriteriaMPC(body)
	log.Println("mpc")
	for _, row := range mpc {
		log.Println(row)
	}
	colSum := ahp.CalculateColSum(mpc)
	log.Printf("colSum: %v", colSum)
	normMpc := ahp.NormalizeMPC(mpc, colSum)
	log.Println("normMpc")
	for _, row := range normMpc {
		log.Println(row)
	}
	criteriaWeight := ahp.CalculateCriteriaWeight(normMpc)
	log.Printf("criteriaWeight: %v", criteriaWeight)

	weightedSum := ahp.CalculateWeightedSum(mpc, criteriaWeight)
	log.Printf("weightedSum: %v", weightedSum)

	lambdaMax := ahp.CalculateLambdaMax(weightedSum, criteriaWeight)
	log.Printf("lambdaMax: %v", lambdaMax)

	consistencyIndex := ahp.ConsistencyIndex(lambdaMax)
	log.Printf("ci: %v", consistencyIndex)

	consistencyRatio := ahp.ConsistencyRatio(consistencyIndex)
	log.Printf("cr: %v", consistencyRatio)

	log.Printf("is consistent: %v", ahp.IsAnswerConsistent(consistencyRatio))

	allSubCriteria := [ahp.TotalCriteria][ahp.TotalSubCriteria]float32{}
	log.Println("allSubCriteria")
	for i := 0; i < ahp.TotalCriteria; i++ {
		subMpc := ahp.BuildSubMPC()
		subColSum := ahp.CalculateSubColSum(subMpc)
		subNormMpc := ahp.NormalizeSubMPC(subMpc, subColSum)
		subCriteriwaWeight := ahp.CalculateSubCriteriaWeight(subNormMpc)
		allSubCriteria[i] = subCriteriwaWeight
		log.Println(allSubCriteria[i])
	}
	log.Println("success create sub criteria")

	schoolIdClaim, err := jwt.GetJwtClaim(r, "school_id")
	if err != nil {
		return errors.New("invalid token")
	}
	schoolId := schoolIdClaim.(string)
	settings, err := repositories.GetQuestionnareSettingRepository().GetQuestionnareSettings(schoolId)
	if err != nil {
		log.Println(err)
		return httputil.ErrInternalServer
	}
	if len(settings) < ahp.TotalAlternative {
		return errors.New("pengaturan kuesioner belum lengkap, seilahkan hubungi admin sekolah anda")
	}

	decisionMatrix, err := ahp.BuildDecisionMatrix(settings, body)
	log.Println("success create decisionMatrix")
	if err != nil {
		log.Println(err)
		return httputil.ErrInternalServer
	}

	alternativeMatrix := [ahp.TotalAlternative][ahp.TotalCriteria]float32{}
	alternativeHpt := []float32{}

	for row := range ahp.TotalAlternative {
		for col := range ahp.TotalCriteria {
			subVecIdx := decisionMatrix[row][col] - 1
			alternativeMatrix[row][col] = allSubCriteria[col][int(subVecIdx)] * criteriaWeight[col]
		}
		alternativeHpt = append(alternativeHpt, ahp.SumRow(alternativeMatrix[row]))
	}

	log.Println("alternative hpt: ", alternativeHpt)
	return nil
}
