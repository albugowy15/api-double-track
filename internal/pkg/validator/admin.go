package validator

import (
	"errors"

	"github.com/albugowy15/api-double-track/internal/pkg/models/user"
)

func ValidateUpdateAdminRequest(data user.UpdateAdminRequest) error {
	if len(data.Username) == 0 {
		return errors.New("username wajib diisi")
	}

	if err := ValidateUsername(data.Username); err != nil {
		return err
	}

	if len(data.Email) > 0 {
		// email validation
		err := ValidateEmail(data.Email)
		if err != nil {
			return err
		}
	}

	if len(data.PhoneNumber) > 0 {
		// phone_number validation
		err := ValidatePhoneNumber(data.PhoneNumber)
		if err != nil {
			return err
		}
	}
	return nil
}
