package controllers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/albugowy15/api-double-track/internal/api/services"
	"github.com/albugowy15/api-double-track/internal/pkg/db"
	"github.com/albugowy15/api-double-track/internal/pkg/models"
	"github.com/albugowy15/api-double-track/internal/pkg/repositories"
	"github.com/albugowy15/api-double-track/internal/pkg/schemas"
	"github.com/albugowy15/api-double-track/internal/pkg/utils/httputil"
	"github.com/albugowy15/api-double-track/internal/pkg/utils/jwt"
	"github.com/albugowy15/api-double-track/internal/pkg/validator"
	"github.com/guregu/null/v5"
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
//	@Success		201				{object}	httputil.MessageJsonResponse
//	@Failure		400				{object}	httputil.ErrorJsonResponse
//	@Failure		500				{object}	httputil.ErrorJsonResponse
//	@Router			/questionnare/settings [post]
func AddQuestionnareSettings(w http.ResponseWriter, r *http.Request) {
	var body models.QuestionnareSetting
	httputil.GetBody(w, r, &body)
	err := validator.ValidateQuestionnareSettings(body)
	if err != nil {
		httputil.SendError(w, err, http.StatusBadRequest)
		return
	}

	schoolIdClaim, _ := jwt.GetJwtClaim(r, "school_id")
	schoolId := schoolIdClaim.(string)
	body.SchoolId = schoolId

	s := repositories.GetQuestionnareSettingRepository()
	err = s.AddQuestionnareSetting(body)
	if err != nil {
		log.Println(err)
		httputil.SendError(w, httputil.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	httputil.SendMessage(w, "berhasil menyimpan setting kuesioner", http.StatusCreated)
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
//	@Success		200				{object}	httputil.DataJsonResponse{data=[]schemas.QuestionResponse}
//	@Failure		400				{object}	httputil.ErrorJsonResponse
//	@Failure		500				{object}	httputil.ErrorJsonResponse
//	@Router			/questionnare/questions [get]
func GetQuestions(w http.ResponseWriter, r *http.Request) {
	questions, err := repositories.GetQuestionRepository().GetQuestions()
	if err != nil {
		httputil.SendError(w, httputil.ErrInternalServer, http.StatusInternalServerError)
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
			httputil.SendError(w, httputil.ErrInternalServer, http.StatusInternalServerError)
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
				httputil.SendError(w, httputil.ErrInternalServer, http.StatusInternalServerError)
				return
			}
		case "COMPARISON":
			item.Options = []string{"9", "7", "5", "3", "1", "1/3", "1/5", "1/7", "1/9"}
			item.MinText = CodeToText[token[0]]
			item.MaxText = CodeToText[token[1]]
		default:
			log.Println("there is no question with category: ", question.Category)
			httputil.SendError(w, httputil.ErrInternalServer, http.StatusInternalServerError)
			return
		}

		questionsRes = append(questionsRes, item)
	}
	httputil.SendData(w, questionsRes, http.StatusOK)
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
//	@Param			body			body		[]schemas.SubmitAnswerRequest	true	"Submit answer request body"
//	@Success		201				{object}	httputil.MessageJsonResponse
//	@Failure		400				{object}	httputil.ErrorJsonResponse
//	@Failure		500				{object}	httputil.ErrorJsonResponse
//	@Router			/questionnare/answers [post]
func SubmitAnswer(w http.ResponseWriter, r *http.Request) {
	var body []models.SubmitAnswerRequest
	httputil.GetBody(w, r, &body)
	schoolIdClaim, _ := jwt.GetJwtClaim(r, "school_id")
	schoolId := schoolIdClaim.(string)
	studentIdClaim, _ := jwt.GetJwtClaim(r, "user_id")
	studentId := studentIdClaim.(string)

	if err := validator.ValidateSubmitAnswer(body); err != nil {
		httputil.SendError(w, err, http.StatusBadRequest)
		return
	}
	if err := validator.ValidateAnswerNumber(body); err != nil {
		httputil.SendError(w, err, http.StatusBadRequest)
		return
	}

	// make sure all questionnare settings set
	missingSettings, err := repositories.GetQuestionnareSettingRepository().GetMissingSettings(schoolId)
	if err != nil {
		httputil.SendError(w, httputil.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	isSettingsComplete := len(missingSettings) == 0
	if !isSettingsComplete {
		httputil.SendError(w, errors.New("kuesioner belum diatur, hubungi admin"), http.StatusBadRequest)
		return
	}
	// make sure student not submit answer twice
	prevAnswer, err := repositories.GetAnswersRepository().GetAnswersByStudentId(studentId)
	if err != nil {
		httputil.SendError(w, httputil.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	isPrevAnswerExist := len(prevAnswer) != 0
	if isPrevAnswerExist {
		httputil.SendError(w, errors.New("kuesioner telah diselesaikan"), http.StatusBadRequest)
		return
	}

	tx, err := db.GetDb().Beginx()
	if err != nil {
		httputil.SendError(w, err, http.StatusBadRequest)
		return
	}

	if err := services.CalculateAHP(r, body, tx); err != nil {
		tx.Rollback()
		re, ok := err.(*services.AHPServiceError)
		if ok {
			httputil.SendError(w, re.Err, re.StatusCode)
			return
		}
		httputil.SendError(w, httputil.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	// start your topsis service here,
	// the function parameter would be same as CalculateAHP services
	// TODO: TOPSIS Service

	// save answer
	answers := []models.Answer{}
	for _, item := range body {
		answer := models.Answer{
			StudentId:  studentId,
			QuestionId: item.Id,
			Answer:     null.StringFrom(item.Answer),
		}
		answers = append(answers, answer)
	}
	err = repositories.GetAnswersRepository().SaveAnswersTx(answers, tx)
	if err != nil {
		tx.Rollback()
		httputil.SendError(w, httputil.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	tx.Commit()
	httputil.SendMessage(w, "berhasil menyimpan kuesioner", http.StatusCreated)
}

// DeleteAnswer godoc
//
//	@Summary		Delete student questionnare answer
//	@Description	Delete student questionnare answer
//	@Tags			Questionnare
//	@Tags			Student
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		201				{object}	httputil.MessageJsonResponse
//	@Failure		400				{object}	httputil.ErrorJsonResponse
//	@Failure		500				{object}	httputil.ErrorJsonResponse
//	@Router			/questionnare/answers [delete]
func DeleteAnswer(w http.ResponseWriter, r *http.Request) {
	studentIdClaim, _ := jwt.GetJwtClaim(r, "user_id")
	studentId := studentIdClaim.(string)
	if err := repositories.GetAnswersRepository().DeleteAnswers(studentId); err != nil {
		httputil.SendError(w, httputil.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	httputil.SendMessage(w, "berhasil menghapus jawaban kuesioner", http.StatusCreated)
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
//	@Success		200				{object}	httputil.DataJsonResponse{data=[]schemas.QuestionnareSettingAlternative}
//	@Failure		400				{object}	httputil.ErrorJsonResponse
//	@Failure		500				{object}	httputil.ErrorJsonResponse
//	@Router			/questionnare/settings/incomplete [get]
func GetIncompleteQuestionnareSettings(w http.ResponseWriter, r *http.Request) {
	schoolIdClaim, _ := jwt.GetJwtClaim(r, "school_id")
	schoolId := schoolIdClaim.(string)
	alternatives, err := repositories.GetQuestionnareSettingRepository().GetMissingSettings(schoolId)
	if err != nil {
		log.Println(err)
		httputil.SendError(w, httputil.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	httputil.SendData(w, alternatives, http.StatusOK)
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
//	@Success		200				{object}	httputil.DataJsonResponse{data=[]schemas.QuestionnareSettingAlternative}
//	@Failure		400				{object}	httputil.ErrorJsonResponse
//	@Failure		500				{object}	httputil.ErrorJsonResponse
//	@Router			/questionnare/settings [get]
func GetQuestionnareSettings(w http.ResponseWriter, r *http.Request) {
	schoolIdClaim, _ := jwt.GetJwtClaim(r, "school_id")
	schoolId := schoolIdClaim.(string)
	settings, err := repositories.GetQuestionnareSettingRepository().GetQuestionnareSettings(schoolId)
	if err != nil {
		log.Println(err)
		httputil.SendError(w, httputil.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	httputil.SendData(w, settings, http.StatusOK)
}

// GetQuesionnareStatus godoc
//
//	@Summary		Get questionnare status
//	@Description	Get questionnare status
//	@Tags			Questionnare
//	@Tags			Student
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		200				{object}	httputil.DataJsonResponse{data=schemas.QuestionnareStatusResponse}
//	@Failure		400				{object}	httputil.ErrorJsonResponse
//	@Failure		500				{object}	httputil.ErrorJsonResponse
//	@Router			/questionnare/status [get]
func GetQuesionnareStatus(w http.ResponseWriter, r *http.Request) {
	schoolIdClaim, _ := jwt.GetJwtClaim(r, "school_id")
	schoolId := schoolIdClaim.(string)
	studentIdClaim, _ := jwt.GetJwtClaim(r, "user_id")
	studentId := studentIdClaim.(string)

	// check is questionnare complete
	answers, err := repositories.GetAnswersRepository().GetAnswersByStudentId(studentId)
	if err != nil {
		httputil.SendError(w, httputil.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	if len(answers) > 0 {
		res := schemas.QuestionnareStatusResponse{
			Status: "COMPLETED",
		}
		httputil.SendData(w, res, http.StatusOK)
		return
	}

	// check is questionnare settings all sets
	settings, err := repositories.GetQuestionnareSettingRepository().GetQuestionnareSettings(schoolId)
	if err != nil {
		httputil.SendError(w, httputil.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	totalAlternative := 7
	if len(settings) == totalAlternative {
		res := schemas.QuestionnareStatusResponse{
			Status: "READY",
		}
		httputil.SendData(w, res, http.StatusOK)
		return
	}

	res := schemas.QuestionnareStatusResponse{
		Status: "NOTREADY",
	}
	httputil.SendData(w, res, http.StatusOK)
}
