package controllers

import (
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/albugowy15/api-double-track/db"
	"github.com/albugowy15/api-double-track/internal/models"
	"github.com/albugowy15/api-double-track/internal/repositories"
	"github.com/albugowy15/api-double-track/internal/services"
	weightmethods "github.com/albugowy15/api-double-track/internal/services/weight_methods"
	"github.com/albugowy15/api-double-track/internal/validator"
	"github.com/albugowy15/api-double-track/pkg/ahp"
	"github.com/albugowy15/api-double-track/pkg/auth"
	"github.com/albugowy15/api-double-track/pkg/httpx"
	"github.com/albugowy15/api-double-track/pkg/schemas"
	"github.com/guregu/null/v5"
)

var CodeToText = map[string]string{
	"JLP": "Jumlah lapangan pekerjaan lebih penting",
	"GAJ": "Gaji lebih penting",
	"PEW": "Peluang wirausaha lebih penting",
	"MIN": "Minat lebih penting",
	"FAS": "Fasilitas pendukung lebih penting",
}

// HandlePostQuestionnareSettings godoc
//
//	@Summary		Add questionnare setting
//	@Description	Add a questionnare setting
//	@Tags			Questionnare
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string								true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			body			body		schemas.QuestionnareSettingRequest	true	"Add questionnare setting request body"
//	@Success		201				{object}	httpx.MessageJsonResponse
//	@Failure		400				{object}	httpx.ErrorJsonResponse
//	@Failure		500				{object}	httpx.ErrorJsonResponse
//	@Router			/questionnare/settings [post]
func HandlePostQuestionnareSettings(w http.ResponseWriter, r *http.Request) {
	var body models.QuestionnareSetting
	if err := httpx.GetBody(r, &body); err != nil {
		httpx.SendError(w, err, http.StatusBadRequest)
	}
	err := validator.ValidateQuestionnareSettings(body)
	if err != nil {
		httpx.SendError(w, err, http.StatusBadRequest)
		return
	}

	schoolIdClaim, _ := auth.GetJwtClaim(r, "school_id")
	schoolId := schoolIdClaim.(string)
	body.SchoolId = schoolId

	err = repositories.AddQuestionnareSetting(body)
	if err != nil {
		log.Println(err)
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	httpx.SendMessage(w, "berhasil menyimpan setting kuesioner", http.StatusCreated)
}

// HandleGetQuestions godoc
//
//	@Summary		Get Questions
//	@Description	Get all available questions
//	@Tags			Questionnare
//	@Tags			Student
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		200				{object}	httpx.DataJsonResponse{data=[]schemas.QuestionResponse}
//	@Failure		400				{object}	httpx.ErrorJsonResponse
//	@Failure		500				{object}	httpx.ErrorJsonResponse
//	@Router			/questionnare/questions [get]
func HandleGetQuestions(w http.ResponseWriter, r *http.Request) {
	questions, err := repositories.GetQuestions()
	if err != nil {
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
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
			httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
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
				httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
				return
			}
		case "COMPARISON":
			item.Options = []string{"9", "7", "5", "3", "1", "1/3", "1/5", "1/7", "1/9"}
			item.MinText = CodeToText[token[0]]
			item.MaxText = CodeToText[token[1]]
		default:
			log.Println("there is no question with category: ", question.Category)
			httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
			return
		}

		questionsRes = append(questionsRes, item)
	}
	httpx.SendData(w, questionsRes, http.StatusOK)
}

// HandlePostAnswers godoc
//
//	@Summary		Submit answer
//	@Description	Submit questionnare answers
//	@Tags			Questionnare
//	@Tags			Student
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Param			body			body		[]schemas.SubmitAnswerRequest	true	"Submit answer request body"
//	@Success		201				{object}	httpx.MessageJsonResponse
//	@Failure		400				{object}	httpx.ErrorJsonResponse
//	@Failure		500				{object}	httpx.ErrorJsonResponse
//	@Router			/questionnare/answers [post]
func HandlePostAnswers(w http.ResponseWriter, r *http.Request) {
	var body []models.SubmitAnswerRequest
	if err := httpx.GetBody(r, &body); err != nil {
		httpx.SendError(w, err, http.StatusBadRequest)
	}
	schoolIdClaim, _ := auth.GetJwtClaim(r, "school_id")
	schoolId := schoolIdClaim.(string)
	studentIdClaim, _ := auth.GetJwtClaim(r, "user_id")
	studentId := studentIdClaim.(string)

	if err := validator.ValidateSubmitAnswer(body); err != nil {
		httpx.SendError(w, err, http.StatusBadRequest)
		return
	}
	if err := validator.ValidateAnswerNumber(body); err != nil {
		httpx.SendError(w, err, http.StatusBadRequest)
		return
	}

	// make sure all questionnare settings set
	missingSettings, err := repositories.GetMissingSettings(schoolId)
	if err != nil {
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	isSettingsComplete := len(missingSettings) == 0
	if !isSettingsComplete {
		httpx.SendError(w, errors.New("kuesioner belum diatur, hubungi admin"), http.StatusBadRequest)
		return
	}
	// make sure student not submit answer twice
	prevAnswer, err := repositories.GetAnswersByStudentId(studentId)
	if err != nil {
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	isPrevAnswerExist := len(prevAnswer) != 0
	if isPrevAnswerExist {
		httpx.SendError(w, errors.New("kuesioner telah diselesaikan"), http.StatusBadRequest)
		return
	}

	tx, err := db.AppDB.Beginx()
	if err != nil {
		httpx.SendError(w, err, http.StatusBadRequest)
		return
	}

	//
	mpc := ahp.BuildCriteriaMPC(body)
	colSum := ahp.CalculateColSum(mpc)
	normMpc := ahp.NormalizeMPC(mpc, colSum)
	criteriaWeight := ahp.CalculateCriteriaWeight(normMpc)

	if err := services.CalculateAHP(r, body, mpc, criteriaWeight, tx); err != nil {
		tx.Rollback()
		re, ok := err.(*services.AHPServiceError)
		if ok {
			httpx.SendError(w, re.Err, re.StatusCode)
			return
		}
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	// start your topsis service here,
	// the function parameter would be same as CalculateAHP services
	// TODO: TOPSIS Service
	if err := weightmethods.CalculateEntropy(r, body); err != nil {
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	if err := services.CalculateTopsis(r, body, tx); err != nil {
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}

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
	err = repositories.SaveAnswersTx(answers, tx)
	if err != nil {
		tx.Rollback()
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}

	tx.Commit()
	httpx.SendMessage(w, "berhasil menyimpan kuesioner", http.StatusCreated)
}

// HandleDeleteAnswer godoc
//
//	@Summary		Delete student questionnare answer
//	@Description	Delete student questionnare answer
//	@Tags			Questionnare
//	@Tags			Student
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string				true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		201				{object}	httpx.MessageJsonResponse
//	@Failure		400				{object}	httpx.ErrorJsonResponse
//	@Failure		500				{object}	httpx.ErrorJsonResponse
//	@Router			/questionnare/answers [delete]
func HandleDeleteAnswer(w http.ResponseWriter, r *http.Request) {
	studentIdClaim, _ := auth.GetJwtClaim(r, "user_id")
	studentId := studentIdClaim.(string)
	if err := repositories.DeleteAnswers(studentId); err != nil {
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	httpx.SendMessage(w, "berhasil menghapus jawaban kuesioner", http.StatusCreated)
}

// HandleGetIncompleteQuestionnareSettings godoc
//
//	@Summary		Get incomplete questionnare settings
//	@Description	Get incomplete questionnare settings
//	@Tags			Questionnare
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		200				{object}	httpx.DataJsonResponse{data=[]schemas.QuestionnareSettingAlternative}
//	@Failure		400				{object}	httpx.ErrorJsonResponse
//	@Failure		500				{object}	httpx.ErrorJsonResponse
//	@Router			/questionnare/settings/incomplete [get]
func HandleGetIncompleteQuestionnareSettings(w http.ResponseWriter, r *http.Request) {
	schoolIdClaim, _ := auth.GetJwtClaim(r, "school_id")
	schoolId := schoolIdClaim.(string)
	alternatives, err := repositories.GetMissingSettings(schoolId)
	if err != nil {
		log.Println(err)
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	httpx.SendData(w, alternatives, http.StatusOK)
}

// HandleGetQuestionnareSettings godoc
//
//	@Summary		Get questionnare settings
//	@Description	Get questionnare settings
//	@Tags			Questionnare
//	@Tags			Admin
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		200				{object}	httpx.DataJsonResponse{data=[]schemas.QuestionnareSettingAlternative}
//	@Failure		400				{object}	httpx.ErrorJsonResponse
//	@Failure		500				{object}	httpx.ErrorJsonResponse
//	@Router			/questionnare/settings [get]
func HandleGetQuestionnareSettings(w http.ResponseWriter, r *http.Request) {
	schoolIdClaim, _ := auth.GetJwtClaim(r, "school_id")
	schoolId := schoolIdClaim.(string)
	settings, err := repositories.GetQuestionnareSettings(schoolId)
	if err != nil {
		log.Println(err)
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	httpx.SendData(w, settings, http.StatusOK)
}

// HandleGetQuesionnareStatus godoc
//
//	@Summary		Get questionnare status
//	@Description	Get questionnare status
//	@Tags			Questionnare
//	@Tags			Student
//	@Accept			json
//	@Produce		json
//	@Param			Authorization	header		string	true	"Insert your access token"	default(Bearer <Add access token here>)
//	@Success		200				{object}	httpx.DataJsonResponse{data=schemas.QuestionnareStatusResponse}
//	@Failure		400				{object}	httpx.ErrorJsonResponse
//	@Failure		500				{object}	httpx.ErrorJsonResponse
//	@Router			/questionnare/status [get]
func HandleGetQuesionnareStatus(w http.ResponseWriter, r *http.Request) {
	schoolIdClaim, _ := auth.GetJwtClaim(r, "school_id")
	schoolId := schoolIdClaim.(string)
	studentIdClaim, _ := auth.GetJwtClaim(r, "user_id")
	studentId := studentIdClaim.(string)

	// check is questionnare complete
	answers, err := repositories.GetAnswersByStudentId(studentId)
	if err != nil {
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	if len(answers) > 0 {
		res := schemas.QuestionnareStatusResponse{
			Status: "COMPLETED",
		}
		httpx.SendData(w, res, http.StatusOK)
		return
	}

	// check is questionnare settings all sets
	settings, err := repositories.GetQuestionnareSettings(schoolId)
	if err != nil {
		httpx.SendError(w, httpx.ErrInternalServer, http.StatusInternalServerError)
		return
	}
	totalAlternative := 7
	if len(settings) == totalAlternative {
		res := schemas.QuestionnareStatusResponse{
			Status: "READY",
		}
		httpx.SendData(w, res, http.StatusOK)
		return
	}

	res := schemas.QuestionnareStatusResponse{
		Status: "NOTREADY",
	}
	httpx.SendData(w, res, http.StatusOK)
}
