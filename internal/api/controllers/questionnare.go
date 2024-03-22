package controllers

import (
	"log"
	"net/http"
	"strings"

	"github.com/albugowy15/api-double-track/internal/pkg/models"
	"github.com/albugowy15/api-double-track/internal/pkg/repositories"
	"github.com/albugowy15/api-double-track/internal/pkg/utils"
	"github.com/albugowy15/api-double-track/internal/pkg/utils/jwt"
	"github.com/albugowy15/api-double-track/internal/pkg/validator"
)

var CodeToText = map[string]string{
	"JLP": "Jumlah lapangan pekerjaan lebih penting",
	"GAJ": "Gaji lebih penting",
	"PEW": "Peluang wirausaha lebih penting",
	"MIN": "Minat lebih penting",
	"FAS": "Fasilitas pendukung lebih penting",
}

// AddQuestionnareSettings godoc
//
//	@Summary		Add questionnare setting
//	@Description	Add a questionnare setting
//	@Tags			Questionnare
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string								true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			body			body		schemas.QuestionnareSettingRequest	true	"Add questionnare setting request body"
//	@Success		201				{object}	schemas.MessageResponse
//	@Failure		400				{object}	utils.ErrorJsonResponse
//	@Failure		500				{object}	utils.ErrorJsonResponse
//	@Router			/questionnare/settings [post]
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

// GetQuestions godoc
//
//	@Summary		Get Questions
//	@Description	Get all available questions
//	@Tags			Questionnare
//	@Tags			Student
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		200				{object}	utils.DataJsonResponse{data=[]schemas.QuestionResponse}
//	@Failure		400				{object}	utils.ErrorJsonResponse
//	@Failure		500				{object}	utils.ErrorJsonResponse
//	@Router			/questionnare/questions [get]
func GetQuestions(w http.ResponseWriter, r *http.Request) {
	questions, err := repositories.GetQuestionRepository().GetQuestions()
	if err != nil {
		utils.SendError(w, "internal server error", http.StatusInternalServerError)
		return
	}
	questionsRes := []models.QuestionResponse{}
	for _, question := range questions {
		item := models.QuestionResponse{
			Id:       question.Id,
			Question: question.Question,
			Type:     "radio",
			Number:   question.Number,
		}
		token := strings.Split(question.Code, "_")
		if len(token) != 2 {
			log.Printf("token length is not 2: got: %d", len(token))
			utils.SendError(w, "internal server error", http.StatusInternalServerError)
			return
		}
		switch question.Category {
		case "PREFERENCE":
			item.Options = []string{"1", "2", "3", "4"}
			switch token[0] {
			case "MIN":
				item.MinText = "Tidak Berminat"
				item.MaxText = "Sangat Berminat"
			case "FAS":
				item.MinText = "Tidak Mendukung"
				item.MaxText = "Sangat Mendukung"
			default:
				log.Printf("unexpected token %s", token[0])
				utils.SendError(w, "internal server error", http.StatusInternalServerError)
				return
			}
		case "COMPARISON":
			item.Options = []string{"9", "7", "5", "3", "1", "1/3", "1/5", "1/7", "1/9"}
			item.MinText = CodeToText[token[0]]
			item.MaxText = CodeToText[token[1]]
		default:
			log.Println("there is no question with category: ", question.Category)
			utils.SendError(w, "internal server error", http.StatusInternalServerError)
			return
		}

		questionsRes = append(questionsRes, item)
	}
	utils.SendJson(w, questionsRes, http.StatusOK)
}

// SubmitAnswer godoc
//
//	@Summary		Submit answer
//	@Description	Submit questionnare answers
//	@Tags			Questionnare
//	@Tags			Student
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			body			body		map[string]string{}	true	"Submit answer request body"
//	@Success		201				{object}	schemas.MessageResponse
//	@Failure		400				{object}	utils.ErrorJsonResponse
//	@Failure		500				{object}	utils.ErrorJsonResponse
//	@Router			/questionnare/answers [post]
func SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	body := map[string]string{}
	utils.GetBody(w, r, &body)
	if err := validator.ValidateSubmitAnswer(body); err != nil {
		utils.SendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	// start process here
	res := models.MessageResponse{
		Message: "berhasil menyimpan jawabah",
	}
	utils.SendJson(w, res, http.StatusCreated)
}

// GetIncompleteQuestionnareSettings godoc
//
//	@Summary		Get incomplete questionnare settings
//	@Description	Get incomplete questionnare settings
//	@Tags			Questionnare
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		200				{object}	utils.DataJsonResponse{data=[]schemas.QuestionnareSettingAlternative}
//	@Failure		400				{object}	utils.ErrorJsonResponse
//	@Failure		500				{object}	utils.ErrorJsonResponse
//	@Router			/questionnare/settings/incomplete [get]
func GetIncompleteQuestionnareSettings(w http.ResponseWriter, r *http.Request) {
	schoolIdClaim, _ := jwt.GetJwtClaim(r, "school_id")
	schoolId := schoolIdClaim.(string)
	alternatives, err := repositories.GetQuestionnareSettingRepository().GetMissingSettings(schoolId)
	if err != nil {
		log.Println(err)
		utils.SendError(w, "internal server error", http.StatusInternalServerError)
		return
	}
	utils.SendJson(w, alternatives, http.StatusOK)
}

// GetQuestionnareSettings godoc
//
//	@Summary		Get questionnare settings
//	@Description	Get questionnare settings
//	@Tags			Questionnare
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		200				{object}	utils.DataJsonResponse{data=[]schemas.QuestionnareSettingAlternative}
//	@Failure		400				{object}	utils.ErrorJsonResponse
//	@Failure		500				{object}	utils.ErrorJsonResponse
//	@Router			/questionnare/settings [get]
func GetQuestionnareSettings(w http.ResponseWriter, r *http.Request) {
	schoolIdClaim, _ := jwt.GetJwtClaim(r, "school_id")
	schoolId := schoolIdClaim.(string)
	settings, err := repositories.GetQuestionnareSettingRepository().GetQuestionnareSettings(schoolId)
	if err != nil {
		log.Println(err)
		utils.SendError(w, "internal server error", http.StatusInternalServerError)
		return
	}
	utils.SendJson(w, settings, http.StatusOK)
}
