package controllers

import (
	"log"
	"net/http"

	"github.com/albugowy15/api-double-track/internal/pkg/models"
	"github.com/albugowy15/api-double-track/internal/pkg/repositories"
	"github.com/albugowy15/api-double-track/internal/pkg/utils"
	"github.com/albugowy15/api-double-track/internal/pkg/utils/jwt"
	"github.com/albugowy15/api-double-track/internal/pkg/validator"
)

func AddQuestionnareSettings(w http.ResponseWriter, r *http.Request) {
	var body models.QuestionnareSetting
	utils.GetBody(w, r, &body)
	err := validator.ValidateQuestionnareSettings(body)
	if err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	schoolIdClaim, _ := jwt.GetJwtClaim(r, "school_id")
	schoolId := schoolIdClaim.(string)
	body.SchoolId = schoolId

	s := repositories.GetQuestionnareSettingRepository()
	err = s.AddQuestionnareSetting(body)
	if err != nil {
		log.Println(err)
		utils.SendError(w, "internal server error", http.StatusInternalServerError)
		return
	}

	res := models.MessageResponse{
		Message: "berhasil menyimpan setting kuesioner",
	}
	utils.SendJson(w, res, http.StatusCreated)
}

func GetQuestions(w http.ResponseWriter, r *http.Request) {
}

func SubmitAnswer(w http.ResponseWriter, r *http.Request) {
}
