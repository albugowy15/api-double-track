package schemas

type MessageResponse struct {
	Message string `json:"message"`
}

type Statistic struct {
	RegisteredStudents       int     `json:"registered_students"`
	QuestionnareCompleted    int     `json:"questionnare_completed"`
	RecommendationAcceptance float32 `json:"recommendation_acceptance"`
	ConsistencyAvg           float32 `json:"consistency_avg"`
}
