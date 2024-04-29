package validator

import (
	"errors"

	"github.com/albugowy15/api-double-track/internal/models"
)

func ValidateExpectationRequest(body models.ExpectationRequest) error {
	// start validation
	if len(body.Expectations) == 0 {
		return errors.New("ekspektasi wajib diisi")
	}
	if len(body.Expectations) != 7 {
		return errors.New("ekspektasi tidak lengkap")
	}
	rankSet := map[int]bool{}
	alternativeIdSet := map[int]bool{}
	for _, expectation := range body.Expectations {
		if expectation.AlternativeId == 0 {
			return errors.New("id alternative tidak valid")
		}
		if expectation.Rank < 1 || expectation.Rank > 7 {
			return errors.New("ranking ekspektasi tidak valid")
		}
		if rankSet[expectation.Rank] {
			return errors.New("terdapat ranking bidang keterampilan yang sama")
		}
		rankSet[expectation.Rank] = true
		if alternativeIdSet[expectation.AlternativeId] {
			return errors.New("terdapat ranking bidang keterampilan yang sama")
		}
		alternativeIdSet[expectation.AlternativeId] = true
	}

	return nil
	// end validation
}
