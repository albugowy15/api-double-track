package models

import "github.com/guregu/null/v5"

type RecommendationResult struct {
	Alternative string      `db:"alternative" json:"alternative"`
	Description null.String `db:"description" json:"description"`
	Score       null.Float  `db:"score" json:"score"`
	Id          int32       `db:"id" json:"id"`
}

type RecommendationResultWithRank struct {
	Alternative string      `json:"alternative"`
	Description null.String `json:"description"`
	Score       null.Float  `json:"score"`
	Rank        int         `json:"rank"`
	Id          int32       `json:"id"`
}

type AhpRecommendation struct {
	Result           []RecommendationResultWithRank `db:"result" json:"result"`
	ConsistencyRatio null.Float                     `db:"consistency_ratio" json:"consistency_ratio"`
}

type TopsisRecommendation struct {
	Result []RecommendationResult `db:"result" json:"result"`
	Id     int32                  `db:"id" json:"id"`
}

type Recommendation struct {
	Topsis TopsisRecommendation `db:"topsis" json:"topsis"`
	Ahp    AhpRecommendation    `db:"ahp" json:"ahp"`
}

type StudentRecommendation struct {
	StudentId        string                 `db:"student_id" json:"student_id"`
	Fullname         string                 `db:"fullname" json:"fullname"`
	Nisn             string                 `db:"nisn" json:"nisn"`
	AhpResults       []RecommendationResult `db:"ahp_results" json:"ahp_results"`
	TopsisResults    []RecommendationResult `db:"topsis_results" json:"topsis_results"`
	ConsistencyRatio null.Float             `db:"consistency_ratio" json:"consistency_ratio"`
}
