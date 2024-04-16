package validator

import (
	"errors"

	"github.com/albugowy15/api-double-track/internal/models/user"
)

func ValidateUpdateAdminRequest(data user.UpdateAdminRequest) error {
	if len(data.Username) == 0 {
		return errors.New("username wajib diisi")
	}

	if err := ValidateUsername(data.Username); err != nil {
		return err
	}

	if len(data.Email) == 0 {
		return errors.New("email wajib diisi")
	}
	if err := ValidateEmail(data.Email); err != nil {
		return err
	}

	if len(data.PhoneNumber) == 0 {
		return errors.New("nomor hp wajib diisi")
	}
	if err := ValidatePhoneNumber(data.PhoneNumber); err != nil {
		return err
	}

	return nil
}
