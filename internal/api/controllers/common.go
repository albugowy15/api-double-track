package controllers

import (
	"log"
	"net/http"

	"github.com/albugowy15/api-double-track/internal/pkg/models"
	"github.com/albugowy15/api-double-track/internal/pkg/repositories"
	"github.com/albugowy15/api-double-track/internal/pkg/utils"
	"github.com/albugowy15/api-double-track/internal/pkg/utils/jwt"
)

func GetStatistics(w http.ResponseWriter, r *http.Request) {
	res := models.Statistic{
		RegisteredStudents:       123,
		QuestionnareCompleted:    450,
		RecommendationAcceptance: 90.34,
		ConsistencyAvg:           92.54,
	}
	utils.SendJson(w, res, http.StatusOK)
}

func GetAlternatives(w http.ResponseWriter, r *http.Request) {
	s := repositories.GetAlternativeRepository()
	alternatives, err := s.GetAlternatives()
	if err != nil {
		utils.SendError(w, "internal server error", http.StatusInternalServerError)
		return
	}
	utils.SendJson(w, alternatives, http.StatusOK)
}

func GetSchool(w http.ResponseWriter, r *http.Request) {
	// get school_id from token
	schoolIdClaim, _ := jwt.GetJwtClaim(r, "school_id")
	schoolId := schoolIdClaim.(string)

	school, err := repositories.GetSchoolRepository().GetSchoolById(schoolId)
	if err != nil {
		log.Printf("err get school: %v", err)
		utils.SendError(w, "internal server error", http.StatusInternalServerError)
		return
	}
	utils.SendJson(w, school, http.StatusOK)
}
