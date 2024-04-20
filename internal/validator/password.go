package validator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/albugowy15/api-double-track/internal/models"
)

func validatePassword(password string, variant string) error {
	passwordLength := len(password)
	if passwordLength < 8 || passwordLength > 16 {
		errVal := fmt.Sprintf("password %s harus terdiri dari 8 hingga 16 karakter", variant)
		return errors.New(errVal)
	}
	if strings.Contains(password, " ") {
		errVal := fmt.Sprintf("password %s tidak boleh terdapat spasi", variant)
		return errors.New(errVal)
	}
	return nil
}

func ValidateChangePassword(req models.ChangePasswordRequest) error {
	if len(req.OldPassword) == 0 {
		return errors.New("password lama wajib diisi")
	}
	if len(req.NewPassword) == 0 {
		return errors.New("password baru wajib diisi")
	}
	if len(req.ConfirmPassword) == 0 {
		return errors.New("konfirmasi password wajib diisi")
	}

	if err := validatePassword(req.OldPassword, "lama"); err != nil {
		return err
	}
	if err := validatePassword(req.NewPassword, "baru"); err != nil {
		return err
	}
	if err := validatePassword(req.ConfirmPassword, "konfirmasi"); err != nil {
		return err
	}

	if req.ConfirmPassword != req.NewPassword {
		return errors.New("konfirmasi password tidak sama dengan password baru")
	}

	return nil
}
