package ahp_test

import (
	"errors"
	"fmt"
	"sort"
	"testing"

	"github.com/albugowy15/api-double-track/internal/models"
	"github.com/albugowy15/api-double-track/pkg/ahp"
	"github.com/guregu/null/v5"
)

var (
	flip        float32 = 1.0 / 9.0
	expectedMpc         = [5][5]float32{
		{1, flip, flip, flip, flip},
		{9, 1, flip, flip, flip},
		{9, 9, 1, flip, flip},
		{9, 9, 9, 1, flip},
		{9, 9, 9, 9, 1},
	}
	expectedSubMpc = [4][4]float32{
		{1, 0.75, 0.5, 0.25},
		{1.3333334, 1, 0.6666667, 0.33333334},
		{2, 1.5, 1, 0.5},
		{4, 3, 2, 1},
	}
	expectedColSum        = ahp.ColSum{37, 28.11111, 19.222221, 10.333333, 1.4444444}
	expectedSubColSum     = ahp.SubColSum{8.333334, 6.25, 4.166667, 2.0833335}
	expectedNormalizedMpc = [5][5]float32{
		{0.027027028, 0.0039525693, 0.005780347, 0.010752688, 0.07692308},
		{0.24324325, 0.035573125, 0.005780347, 0.010752688, 0.07692308},
		{0.24324325, 0.3201581, 0.052023124, 0.010752688, 0.07692308},
		{0.24324325, 0.3201581, 0.4682081, 0.0967742, 0.07692308},
		{0.24324325, 0.3201581, 0.4682081, 0.87096775, 0.6923077},
	}
	expectedSubNormalizedMpc = [4][4]float32{
		{0.11999999, 0.12, 0.11999999, 0.11999999},
		{0.16, 0.16, 0.16, 0.16},
		{0.23999998, 0.24, 0.23999998, 0.23999998},
		{0.47999996, 0.48, 0.47999996, 0.47999996},
	}
	expectedCriteriaWeight                                             = ahp.CriteriaWeight{0.02488714, 0.074454494, 0.14062004, 0.24106136, 0.518977}
	expectedSubCriteriaWeight                                          = ahp.SubCriteriaWeight{0.11999999, 0.16, 0.23999998, 0.47999996}
	expectedWeightedSum                                                = ahp.WeightedSum{0.13323301, 0.3985119, 1.1191434, 2.4583807, 4.8481846}
	expectedLambdaMax          float32                                 = 7.640901
	expectedConsistencyIndex   float32                                 = 0.6602253
	expectedConsistencyRatio   float32                                 = 0.58948684
	sampleQuestionnareSettings []models.QuestionnareSettingAlternative = []models.QuestionnareSettingAlternative{
		{Alternative: "Multimedia", Id: 1, TotalOpenJobs: null.NewInt16(3, true), EntrepreneurshipOpportunity: null.Int16From(1), Salary: null.Int16From(2)},
		{Alternative: "Tata Boga", Id: 2, TotalOpenJobs: null.NewInt16(1, true), EntrepreneurshipOpportunity: null.Int16From(4), Salary: null.Int16From(1)},
		{Alternative: "Tata Kecantikan", Id: 3, TotalOpenJobs: null.NewInt16(4, true), EntrepreneurshipOpportunity: null.Int16From(1), Salary: null.Int16From(3)},
		{Alternative: "Tata Busana", Id: 4, TotalOpenJobs: null.NewInt16(2, true), EntrepreneurshipOpportunity: null.Int16From(3), Salary: null.Int16From(4)},
		{Alternative: "Teknik Listrik", Id: 5, TotalOpenJobs: null.NewInt16(1, true), EntrepreneurshipOpportunity: null.Int16From(4), Salary: null.Int16From(2)},
		{Alternative: "Teknik Elektro", Id: 6, TotalOpenJobs: null.NewInt16(4, true), EntrepreneurshipOpportunity: null.Int16From(1), Salary: null.Int16From(3)},
		{Alternative: "Teknik Kendaraan Ringan/Motor", Id: 7, TotalOpenJobs: null.NewInt16(3, true), EntrepreneurshipOpportunity: null.Int16From(2), Salary: null.Int16From(2)},
	}
	expectedDecisionMatrix = ahp.DecisionMatrix{
		{3, 2, 1, 2, 2},
		{4, 3, 1, 2, 2},
		{1, 2, 4, 2, 2},
		{2, 4, 3, 2, 2},
		{1, 1, 4, 2, 2},
		{4, 3, 1, 2, 2},
		{3, 2, 2, 2, 2},
	}
	expectedAlternativeHpt = []float32{0.15636617, 0.16829544, 0.20400292, 0.19507504, 0.20102474, 0.16829544, 0.16199097}
)

func TestAhp(t *testing.T) {
	subMpc := ahp.BuildSubMPC()
	if subMpc != expectedSubMpc {
		t.Error("subMpc not match")
	}
	subColSum := ahp.CalculateSubColSum(subMpc)
	if subColSum != expectedSubColSum {
		t.Error("subColSum not match")
	}
	subNormalizedMpc := ahp.NormalizeSubMPC(subMpc, subColSum)
	if subNormalizedMpc != expectedSubNormalizedMpc {
		t.Error("subNormalizedMpc not match")
	}
	subCriteriaWeight := ahp.CalculateSubCriteriaWeight(subNormalizedMpc)
	if subCriteriaWeight != expectedSubCriteriaWeight {
		t.Error("subCriteriaWeight not match")
	}

	body := buildBody()
	mpc := ahp.BuildCriteriaMPC(body)
	if mpc != expectedMpc {
		t.Error("mpc not match")
	}
	if mpc != expectedMpc {
		t.Error("mpc not match")
	}
	colSum := ahp.CalculateColSum(mpc)
	if colSum != expectedColSum {
		t.Error("colSum not match")
	}
	normalizedMpc := ahp.NormalizeMPC(mpc, colSum)
	if normalizedMpc != expectedNormalizedMpc {
		t.Error("normalizedMpc not match")
	}
	criteriaWeight := ahp.CalculateCriteriaWeight(normalizedMpc)
	if criteriaWeight != expectedCriteriaWeight {
		t.Error("criteriaWeight not match")
	}
	weightedSum := ahp.CalculateWeightedSum(mpc, criteriaWeight)
	if weightedSum != expectedWeightedSum {
		t.Error("weightedSum not match")
	}
	lambdaMax := ahp.CalculateLambdaMax(weightedSum, criteriaWeight)
	if lambdaMax != expectedLambdaMax {
		t.Error("lambdaMax not match")
	}
	consistencyIndex := ahp.ConsistencyIndex(lambdaMax)
	if consistencyIndex != expectedConsistencyIndex {
		t.Error("consistencyIndex not match")
	}
	consistencyRatio := ahp.ConsistencyRatio(consistencyIndex)
	if consistencyRatio != expectedConsistencyRatio {
		t.Error("consistencyRatio not match")
	}
	isAnswerConsistent := ahp.IsAnswerConsistent(consistencyRatio)
	if isAnswerConsistent {
		t.Error("expected isAnswerConsistent to be true")
	}
	decisionMatrix, err := ahp.BuildDecisionMatrix(sampleQuestionnareSettings, body)
	if err != nil {
		t.Errorf("got error: %v", err)
	}
	if decisionMatrix != expectedDecisionMatrix {
		t.Error("decisionMatrix not match")
		t.Logf("decisionMatrix: %v\n", decisionMatrix)
	}

	allSubCriteriaWeights := initSubCriteriaWeights()
	alternativeHpt := ahp.CalculateAlternativeHpt(decisionMatrix, allSubCriteriaWeights, criteriaWeight)
	if alternativeHpt == nil {
		t.Error("unexpected alternativeHpt to be nil")
	}
	for idx := range alternativeHpt {
		if alternativeHpt[idx] != expectedAlternativeHpt[idx] {
			t.Error("alternativeHpt not match")
		}
	}

	// cover if alternative name not found
	prevCorrectAlternative := sampleQuestionnareSettings[0].Alternative
	sampleQuestionnareSettings[0].Alternative = "Not Valid Alternative"
	_, err = ahp.BuildDecisionMatrix(sampleQuestionnareSettings, body)
	if err == nil {
		t.Error("expected to return an error")
	}
	message := fmt.Sprintf("no matching row for alternative: %s", "Not Valid Alternative")
	expectedErr := errors.New(message)
	if err.Error() != expectedErr.Error() {
		t.Errorf("expected error to be %v, got %v instead", expectedErr, err)
	}
	sampleQuestionnareSettings[0].Alternative = prevCorrectAlternative

	// cover if total open jobs setting null
	prevCorrectTotalOpenJobs := sampleQuestionnareSettings[0].TotalOpenJobs
	sampleQuestionnareSettings[0].TotalOpenJobs = null.NewInt16(0, false)
	_, err = ahp.BuildDecisionMatrix(sampleQuestionnareSettings, body)
	if err == nil {
		t.Error("expected to return an error")
	}
	expectedErr = errors.New("total open jobs setting is null")
	if err.Error() != expectedErr.Error() {
		t.Errorf("expected error to be %v, got %v instead", expectedErr, err)
	}
	sampleQuestionnareSettings[0].TotalOpenJobs = prevCorrectTotalOpenJobs

	// cover if salary setting null
	prevCorrectSalary := sampleQuestionnareSettings[0].Salary
	sampleQuestionnareSettings[0].Salary = null.NewInt16(0, false)
	_, err = ahp.BuildDecisionMatrix(sampleQuestionnareSettings, body)
	if err == nil {
		t.Error("expected to return an error")
	}
	expectedErr = errors.New("salary setting is null")
	if err.Error() != expectedErr.Error() {
		t.Errorf("expected error to be %v, got %v instead", expectedErr, err)
	}
	sampleQuestionnareSettings[0].Salary = prevCorrectSalary

	// cover if salary setting null
	prevCorrectEntrOppr := sampleQuestionnareSettings[0].EntrepreneurshipOpportunity
	sampleQuestionnareSettings[0].EntrepreneurshipOpportunity = null.NewInt16(0, false)
	_, err = ahp.BuildDecisionMatrix(sampleQuestionnareSettings, body)
	if err == nil {
		t.Error("expected to return an error")
	}
	expectedErr = errors.New("entrepreneurship opportunity setting is null")
	if err.Error() != expectedErr.Error() {
		t.Errorf("expected error to be %v, got %v instead", expectedErr, err)
	}
	sampleQuestionnareSettings[0].EntrepreneurshipOpportunity = prevCorrectEntrOppr
}

func buildBody() []models.SubmitAnswerRequest {
	body := []models.SubmitAnswerRequest{}
	for i := 1; i <= 14; i++ {
		item := models.SubmitAnswerRequest{
			Id:     i,
			Number: i,
			Answer: "2",
		}
		body = append(body, item)
	}
	for i := 15; i <= 24; i++ {
		item := models.SubmitAnswerRequest{
			Id:     i,
			Number: i,
			Answer: "9",
		}
		body = append(body, item)
	}

	sort.Slice(body, func(i, j int) bool {
		return body[i].Number < body[j].Number
	})
	return body
}

func initSubCriteriaWeights() [ahp.TotalCriteria][ahp.TotalSubCriteria]float32 {
	allSubCriteriaWeights := [ahp.TotalCriteria][ahp.TotalSubCriteria]float32{}
	for i := 0; i < ahp.TotalCriteria; i++ {
		subMpc := ahp.BuildSubMPC()
		subColSum := ahp.CalculateSubColSum(subMpc)
		subNormMpc := ahp.NormalizeSubMPC(subMpc, subColSum)
		subCriteriwaWeight := ahp.CalculateSubCriteriaWeight(subNormMpc)
		allSubCriteriaWeights[i] = subCriteriwaWeight
	}
	return allSubCriteriaWeights
}
