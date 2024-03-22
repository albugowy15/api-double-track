package validator

import (
	"errors"
	"strconv"
)

var (
	ErrAnswerLen     = errors.New("semua pertanyaan wajib diisi")
	ErrInvalidId     = errors.New("id pertanyaan tidak valid")
	ErrInvalidAnswer = errors.New("terdapat jawaban yang tidak valid. periksa lagi jawaban anda")
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

func ValidateSubmitAnswer(body map[string]string) error {
	if len(body) != 24 {
		return ErrAnswerLen
	}
	for key, value := range body {
		questionId, err := strconv.Atoi(key)
		if err != nil || questionId == 0 {
			return ErrInvalidId
		}
		if questionId >= 1 && questionId <= 14 {
			answerNum, err := strconv.Atoi(value)
			if err != nil {
				return ErrInvalidAnswer
			}
			_, ok := ValidPrefAnswer[answerNum]
			if !ok {
				return ErrInvalidAnswer
			}
		} else if questionId >= 15 && questionId <= 24 {
			_, ok := ValidCompAnser[value]
			if !ok {
				return ErrInvalidAnswer
			}
		}
	}
	return nil
}
