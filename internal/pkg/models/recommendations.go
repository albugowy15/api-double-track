package models

import "github.com/guregu/null/v5"

type RecommendationResult struct {
	Alternative string     `db:"alternative" json:"alternative"`
	Score       null.Float `db:"score" json:"score"`
	Id          int32      `db:"id" json:"id"`
}

type AhpRecommendation struct {
	Result           RecommendationResult `db:"result" json:"result"`
	ConsistencyIndex null.Float           `db:"consistency_index" json:"consistency_index"`
	Id               int32                `db:"id" json:"id"`
}

type TopsisRecommendation struct {
	Result RecommendationResult `db:"result" json:"result"`
	Id     int32                `db:"id" json:"id"`
}

type Recommendation struct {
	StudentId       string               `db:"student_id" json:"student_id"`
	StudentFullname string               `db:"student_fullname" json:"student_fullname"`
	Consistency     string               `db:"consistency" json:"consistency"`
	Topsis          TopsisRecommendation `db:"topsis" json:"topsis"`
	Ahp             AhpRecommendation    `db:"ahp" json:"ahp"`
}
