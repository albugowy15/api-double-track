package utils_test

import (
	"testing"

	"github.com/albugowy15/api-double-track/internal/models"
	"github.com/albugowy15/api-double-track/internal/utils"
	"github.com/guregu/null/v5"
)

func TestMakeRecommendationRanks(t *testing.T) {
	results := []models.RecommendationResult{
		{
			Id:          1,
			Alternative: "Multimedia",
			Score:       null.FloatFrom(0.3456),
			Description: null.StringFrom("Alternative Multimedia"),
		},
		{
			Id:          2,
			Alternative: "Tata Kecantikan",
			Score:       null.FloatFrom(0.3456),
			Description: null.StringFrom("Tata Kecantikan"),
		},
		{
			Id:          3,
			Alternative: "Tata Boga",
			Score:       null.FloatFrom(0.2345),
			Description: null.StringFrom("Tata Boga"),
		},
		{
			Id:          4,
			Alternative: "Tata Busana",
			Score:       null.FloatFrom(0.2345),
			Description: null.StringFrom("Tata Busana"),
		},
		{
			Id:          5,
			Alternative: "Teknik Listrik",
			Score:       null.FloatFrom(0.2345),
			Description: null.StringFrom("Teknik Listrik"),
		},
		{
			Id:          6,
			Alternative: "Teknik Elektro",
			Score:       null.FloatFrom(0.1234),
			Description: null.StringFrom("Alternative Teknik Elektro"),
		},
	}

	ranks, err := utils.MakeRecommendationRanks(results)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ranks == nil {
		t.Fatal("unexpected ranks nil")
	}

	var expectedRanks []models.RecommendationResultWithRank
	for _, res := range results {
		rank := 1
		if res.Id >= 3 && res.Id <= 5 {
			rank = 2
		} else if res.Id == 6 {
			rank = 3
		}
		expectedRank := models.RecommendationResultWithRank{
			Id:          res.Id,
			Alternative: res.Alternative,
			Score:       res.Score,
			Description: res.Description,
			Rank:        rank,
		}
		expectedRanks = append(expectedRanks, expectedRank)
	}

	for idx := range ranks {
		if ranks[idx] != expectedRanks[idx] {
			t.Fatalf("rank not match, got: %v", ranks[idx])
		}
	}
}
