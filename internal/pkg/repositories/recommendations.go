package repositories

type RecommendationRepository struct{}

var recommendationRepository *RecommendationRepository

func GetRecommendationRepository() *RecommendationRepository {
	if recommendationRepository == nil {
		recommendationRepository = &RecommendationRepository{}
	}
	return recommendationRepository
}

func (r *RecommendationRepository) GetRecommendations(schoolId string) error {
	return nil
}
