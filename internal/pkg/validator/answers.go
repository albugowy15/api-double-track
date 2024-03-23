package validator

import (
	"errors"
	"strconv"

	"github.com/albugowy15/api-double-track/internal/pkg/models"
	"github.com/albugowy15/api-double-track/internal/pkg/repositories"
)

var (
	ErrAnswerLen        = errors.New("semua pertanyaan wajib diisi")
	ErrInvalidId        = errors.New("id pertanyaan tidak valid")
	ErrInvalidAnswer    = errors.New("terdapat jawaban yang tidak valid. periksa lagi jawaban anda")
	ErrDuplicateId      = errors.New("terdapat pertanyaan yang dijawab lebih dari sekali")
	ErrIdNumbetNotMatch = errors.New("id dan nomor pertanyaan tidak cocok")
)

var ValidCompAnser = map[string]bool{
	"9":   true,
	"8":   true,
	"7":   true,
	"5":   true,
	"3":   true,
	"1":   true,
	"1/3": true,
	"1/5": true,
	"1/7": true,
	"1/9": true,
}

var ValidPrefAnswer = map[int]bool{
	1: true,
	2: true,
	3: true,
	4: true,
}

func ValidateSubmitAnswer(body []models.SubmitAnswerRequest) error {
	if len(body) != 24 {
		return ErrAnswerLen
	}
	QuestionIdSet := map[int]bool{}
	for _, item := range body {
		questionId := item.Id
		isExist, ok := QuestionIdSet[questionId]
		if ok || isExist {
			return ErrDuplicateId
		}
		QuestionIdSet[questionId] = true
		if questionId >= 1 && questionId <= 14 {
			answerNum, err := strconv.Atoi(item.Answer)
			if err != nil {
				return ErrInvalidAnswer
			}
			_, ok := ValidPrefAnswer[answerNum]
			if !ok {
				return ErrInvalidAnswer
			}
		} else if questionId >= 15 && questionId <= 24 {
			_, ok := ValidCompAnser[item.Answer]
			if !ok {
				return ErrInvalidAnswer
			}
		} else {
			return ErrInvalidId
		}
	}
	return nil
}

func ValidateAnswerNumber(body []models.SubmitAnswerRequest) error {
	// db query return result as map id to number
	questions, err := repositories.GetQuestionRepository().GetQuestions()
	if err != nil {
		return err
	}
	idToNumberMap := map[int]int{}
	for _, question := range questions {
		idToNumberMap[question.Id] = question.Number
	}

	for _, item := range body {
		questionNum, ok := idToNumberMap[item.Id]
		if !ok {
			return ErrInvalidId
		}
		isNumberMatch := questionNum == item.Number
		if !isNumberMatch {
			return ErrIdNumbetNotMatch
		}
	}
	return nil
}
