package services

import (
	"errors"
	"log"
	"net/http"

	"github.com/albugowy15/api-double-track/internal/pkg/ahp"
	"github.com/albugowy15/api-double-track/internal/pkg/db"
	"github.com/albugowy15/api-double-track/internal/pkg/models"
	"github.com/albugowy15/api-double-track/internal/pkg/repositories"
	"github.com/albugowy15/api-double-track/internal/pkg/utils/httputil"
	"github.com/albugowy15/api-double-track/internal/pkg/utils/jwt"
	"github.com/guregu/null/v5"
)

type AHPServiceError struct {
	Err        error
	StatusCode int
}

func (e *AHPServiceError) Error() string {
	return e.Err.Error()
}

func CalculateAHP(r *http.Request, body []models.SubmitAnswerRequest) error {
	mpc := ahp.BuildCriteriaMPC(body)
	colSum := ahp.CalculateColSum(mpc)
	normMpc := ahp.NormalizeMPC(mpc, colSum)
	criteriaWeight := ahp.CalculateCriteriaWeight(normMpc)
	weightedSum := ahp.CalculateWeightedSum(mpc, criteriaWeight)
	lambdaMax := ahp.CalculateLambdaMax(weightedSum, criteriaWeight)
	consistencyIndex := ahp.ConsistencyIndex(lambdaMax)
	consistencyRatio := ahp.ConsistencyRatio(consistencyIndex)
	log.Printf("is consistent: %v", ahp.IsAnswerConsistent(consistencyRatio))
	allSubCriteria := [ahp.TotalCriteria][ahp.TotalSubCriteria]float32{}
	for i := 0; i < ahp.TotalCriteria; i++ {
		subMpc := ahp.BuildSubMPC()
		subColSum := ahp.CalculateSubColSum(subMpc)
		subNormMpc := ahp.NormalizeSubMPC(subMpc, subColSum)
		subCriteriwaWeight := ahp.CalculateSubCriteriaWeight(subNormMpc)
		allSubCriteria[i] = subCriteriwaWeight
	}

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
			alternativeMatrix[row][col] = allSubCriteria[col][int(subVecIdx)] * criteriaWeight[col]
		}
		alternativeHpt = append(alternativeHpt, ahp.SumRow(alternativeMatrix[row]))
	}

	// save ahp
	studentIdClaim, _ := jwt.GetJwtClaim(r, "user_id")
	studentId := studentIdClaim.(string)
	tx, err := db.GetDb().Beginx()
	if err != nil {
		log.Println(err)
		return &AHPServiceError{
			StatusCode: http.StatusInternalServerError,
			Err:        httputil.ErrInternalServer,
		}
	}
	a := repositories.GetAHPRepository()
	insertedId, err := a.SaveAHPTx(models.AHP{StudentId: studentId, ConsistencyIndex: consistencyIndex}, tx)
	if err != nil {
		tx.Rollback()
		return &AHPServiceError{
			StatusCode: http.StatusInternalServerError,
			Err:        httputil.ErrInternalServer,
		}
	}

	alt := repositories.GetAlternativeRepository()
	alternatives, err := alt.GetAlternatives()
	if err != nil {
		tx.Rollback()
		return &AHPServiceError{
			StatusCode: http.StatusInternalServerError,
			Err:        httputil.ErrInternalServer,
		}
	}

	ahpAlternatives := []models.AHPToAlternatives{}
	for _, alternative := range alternatives {
		hptIdx, ok := ahp.AlternativeToRow[alternative.Alternative]
		if !ok {
			tx.Rollback()
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
		tx.Rollback()
		return &AHPServiceError{
			StatusCode: http.StatusInternalServerError,
			Err:        httputil.ErrInternalServer,
		}
	}

	answers := []models.Answer{}
	for _, item := range body {
		answer := models.Answer{
			StudentId:  studentId,
			QuestionId: item.Id,
			Answer:     null.StringFrom(item.Answer),
		}
		answers = append(answers, answer)
	}
	err = repositories.GetAnswersRepository().SaveAnswersTx(answers, tx)
	if err != nil {
		tx.Rollback()
		return &AHPServiceError{
			StatusCode: http.StatusInternalServerError,
			Err:        httputil.ErrInternalServer,
		}
	}

	err = tx.Commit()
	if err != nil {
		log.Println("err commit transactions:", err)
		return &AHPServiceError{
			StatusCode: http.StatusInternalServerError,
			Err:        httputil.ErrInternalServer,
		}
	}

	return nil
}
