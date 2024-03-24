package ahp_test

import (
	"fmt"
	"sort"
	"testing"

	"github.com/albugowy15/api-double-track/internal/pkg/ahp"
	"github.com/albugowy15/api-double-track/internal/pkg/models"
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
	expectedColSum = ahp.ColSum{37, 28.11111, 19.222221, 10.333333, 1.4444444}
)

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

	fmt.Println(body)
	sort.Slice(body, func(i, j int) bool {
		return body[i].Number < body[j].Number
	})
	return body
}

func TestBuildCriteriaMPC(t *testing.T) {
	body := buildBody()
	mpc := ahp.BuildCriteriaMPC(body)

	if mpc != expectedMpc {
		t.Error("mpc not match")
	}
}

func TestCalculateColSum(t *testing.T) {
	colSum := ahp.CalculateColSum(expectedMpc)
	if colSum != expectedColSum {
		t.Error("colSum not match")
	}
}
