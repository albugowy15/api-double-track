package utils

import (
	"errors"

	"github.com/albugowy15/api-double-track/internal/models"
)

func MakeRecommendationRanks(data []models.RecommendationResult) ([]models.RecommendationResultWithRank, error) {
	var ranks []models.RecommendationResultWithRank

	if len(data) == 0 {
		return nil, errors.New("recommendation data is empty")
	}
	if !data[0].Score.Valid {
		return nil, errors.New("recommendation score is null")
	}
	maxScore := data[0].Score.Float64
	currRank := 1

	for _, item := range data {
		currScore := item.Score
		if !currScore.Valid {
			return nil, errors.New("recommendation score is null")
		}
		if currScore.Float64 < maxScore {
			maxScore = currScore.Float64
			currRank += 1
		}
		rec := models.RecommendationResultWithRank{
			Rank:        currRank,
			Score:       currScore,
			Description: item.Description,
			Id:          item.Id,
			Alternative: item.Alternative,
		}
		ranks = append(ranks, rec)
	}

	return ranks, nil
}
