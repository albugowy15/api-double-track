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

func SubmitAnswer(w http.ResponseWriter, r *http.Request) {
}
