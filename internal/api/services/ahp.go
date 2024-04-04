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
	"github.com/jmoiron/sqlx"
)

type AHPServiceError struct {
	Err        error
	StatusCode int
}

func (e *AHPServiceError) Error() string {
	return e.Err.Error()
}

// this allSubCriteriaWeights will be initialized inside main function
var allSubCriteriaWeights = [ahp.TotalCriteria][ahp.TotalSubCriteria]float32{}

func InitSubCriteriaWeights() {
	for i := 0; i < ahp.TotalCriteria; i++ {
		subMpc := ahp.BuildSubMPC()
		subColSum := ahp.CalculateSubColSum(subMpc)
		subNormMpc := ahp.NormalizeSubMPC(subMpc, subColSum)
		subCriteriwaWeight := ahp.CalculateSubCriteriaWeight(subNormMpc)
		allSubCriteriaWeights[i] = subCriteriwaWeight
	}
}

func CalculateAHP(r *http.Request, body []models.SubmitAnswerRequest, tx *sqlx.Tx) error {
	mpc := ahp.BuildCriteriaMPC(body)
	colSum := ahp.CalculateColSum(mpc)
	normMpc := ahp.NormalizeMPC(mpc, colSum)
	criteriaWeight := ahp.CalculateCriteriaWeight(normMpc)
	weightedSum := ahp.CalculateWeightedSum(mpc, criteriaWeight)
	lambdaMax := ahp.CalculateLambdaMax(weightedSum, criteriaWeight)
	consistencyIndex := ahp.ConsistencyIndex(lambdaMax)
	consistencyRatio := ahp.ConsistencyRatio(consistencyIndex)
	log.Printf("is consistent: %v", ahp.IsAnswerConsistent(consistencyRatio))

	schoolIdClaim, err := jwt.GetJwtClaim(r, "school_id")
	if err != nil {
		return &AHPServiceError{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("invalid token"),
		}
	}
	schoolId := schoolIdClaim.(string)
	settings, err := repositories.GetQuestionnareSettingRepository().GetQuestionnareSettings(schoolId)
	if err != nil {
		log.Println(err)
		return &AHPServiceError{
			StatusCode: http.StatusInternalServerError,
			Err:        httputil.ErrInternalServer,
		}
	}
	if len(settings) < ahp.TotalAlternative {
		return &AHPServiceError{
			StatusCode: http.StatusBadRequest,
			Err:        errors.New("pengaturan kuesioner belum lengkap, seilahkan hubungi admin sekolah anda"),
		}
	}

	decisionMatrix, err := ahp.BuildDecisionMatrix(settings, body)
	if err != nil {
		log.Println(err)
		return &AHPServiceError{
			StatusCode: http.StatusInternalServerError,
			Err:        httputil.ErrInternalServer,
		}
	}

	alternativeMatrix := [ahp.TotalAlternative][ahp.TotalCriteria]float32{}
	alternativeHpt := []float32{}
	for row := range ahp.TotalAlternative {
		for col := range ahp.TotalCriteria {
			subVecIdx := decisionMatrix[row][col] - 1
			alternativeMatrix[row][col] = allSubCriteriaWeights[col][int(subVecIdx)] * criteriaWeight[col]
		}
		alternativeHpt = append(alternativeHpt, ahp.SumRow(alternativeMatrix[row]))
	}

	// save ahp
	studentIdClaim, _ := jwt.GetJwtClaim(r, "user_id")
	studentId := studentIdClaim.(string)
	a := repositories.GetAHPRepository()
	insertedId, err := a.SaveAHPTx(models.AHP{StudentId: studentId, ConsistencyRatio: consistencyRatio}, tx)
	if err != nil {
		return &AHPServiceError{
			StatusCode: http.StatusInternalServerError,
			Err:        httputil.ErrInternalServer,
		}
	}

	alt := repositories.GetAlternativeRepository()
	alternatives, err := alt.GetAlternatives()
	if err != nil {
		return &AHPServiceError{
			StatusCode: http.StatusInternalServerError,
			Err:        httputil.ErrInternalServer,
		}
	}

	ahpAlternatives := []models.AHPToAlternatives{}
	for _, alternative := range alternatives {
		hptIdx, ok := ahp.AlternativeToRow[alternative.Alternative]
		if !ok {
			log.Printf("err %s is not valid index for AlternativeToRow\n", alternative.Alternative)
			return &AHPServiceError{
				StatusCode: http.StatusInternalServerError,
				Err:        httputil.ErrInternalServer,
			}
		}
		score := alternativeHpt[hptIdx]
		ahpAlternative := models.AHPToAlternatives{
			Score:         score,
			AlternativeId: alternative.Id,
			AhpId:         insertedId,
		}
		ahpAlternatives = append(ahpAlternatives, ahpAlternative)
	}
	err = a.SaveAHPAlternativesTx(ahpAlternatives, tx)
	if err != nil {
		return &AHPServiceError{
			StatusCode: http.StatusInternalServerError,
			Err:        httputil.ErrInternalServer,
		}
	}

	return nil
}
