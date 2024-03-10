package models

type MessageResponse struct {
	Message string `json:"message"`
}

type Statistic struct {
	RegisteredStudents       int     `db:"registered_students" json:"registered_students"`
	QuestionnareCompleted    int     `db:"questionnare_completed" json:"questionnare_completed"`
	RecommendationAcceptance float32 `db:"recommendation_acceptance" json:"recommendation_acceptance"`
	ConsistencyAvg           float32 `db:"consistency_avg" json:"consistency_avg"`
}
