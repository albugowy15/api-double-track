package schemas

type RecommendationResult struct {
	Alternative string  `json:"alternative"`
	Score       float32 `json:"score"`
	Id          int32   `json:"id"`
}

type AhpRecommendation struct {
	Result           RecommendationResult `json:"result"`
	ConsistencyIndex float32              `json:"consistency_index"`
	Id               int32                `json:"id"`
}

type TopsisRecommendation struct {
	Result RecommendationResult `json:"result"`
	Id     int32                `json:"id"`
}

type Recommendation struct {
	StudentId       string               `json:"student_id"`
	StudentFullname string               `json:"student_fullname"`
	Consistency     string               `json:"consistency"`
	Topsis          TopsisRecommendation `json:"topsis"`
	Ahp             AhpRecommendation    `json:"ahp"`
}
