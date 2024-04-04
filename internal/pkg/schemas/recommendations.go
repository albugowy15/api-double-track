package schemas

type RecommendationResult struct {
	Alternative string  `json:"alternative"`
	Description string  `json:"description"`
	Score       float32 `json:"score"`
	Id          int32   `json:"id"`
}

type AhpRecommendation struct {
	Result           RecommendationResult `json:"result"`
	ConsistencyRatio float32              `json:"consistency_ratio"`
}

type TopsisRecommendation struct {
	Result RecommendationResult `json:"result"`
	Id     int32                `json:"id"`
}

type Recommendation struct {
	Ahp    AhpRecommendation    `json:"ahp"`
	Topsis TopsisRecommendation `json:"topsis"`
}

type StudentRecommendation struct {
	Fullname         string                 `json:"fullname"`
	Nisn             string                 `json:"nisn"`
	AhpResults       []RecommendationResult `json:"ahp_results"`
	TopsisResults    []RecommendationResult `json:"topsis_results"`
	Id               int32                  `json:"id"`
	ConsistencyRatio float32                `json:"consistency_ratio"`
}

type DeleteRecommendationRequest struct {
	StudentId string `json:"student_id"`
}
