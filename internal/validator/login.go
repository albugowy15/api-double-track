package validator

import (
	"errors"
	"regexp"

	"github.com/albugowy15/api-double-track/internal/models"
)

func ValidateUsername(username string) error {
	usernameLen := len(username)
	if usernameLen < 6 || usernameLen > 30 {
		return errors.New("username harus terdiri dari 6 hingga 30 karakter")
	}

	isWhitespacePresent := regexp.MustCompile(`\s`).MatchString(username)
	if isWhitespacePresent {
		return ErrUsernameWhitespace
	}
	pattern := "^[a-zA-Z0-9]+$"
	re := regexp.MustCompile(pattern)
	if !re.MatchString(username) {
		return errors.New("username hanya boleh terdiri dari huruf alphabet atau angka tanpa spasi")
	}
	return nil
}

func ValidatePassword(password string) error {
	passwordLength := len(password)
	if passwordLength == 0 {
		return errors.New("password wajib diisi")
	}
	if len(password) < 6 {
		return ErrUsernameLength
	}
	isWhitespacePresent := regexp.MustCompile(`\s`).MatchString(password)
	if isWhitespacePresent {
		return ErrPasswordWhitespace
	}
	return nil
}

func ValidateLoginType(loginType string) error {
	if loginType != "admin" && loginType != "student" {
		return errors.New("type login tidak valid")
	}
	return nil
}

func ValidateLoginRequest(body models.LoginRequest) error {
	if err := ValidateUsername(body.Username); err != nil {
		return err
	}
	if err := ValidatePassword(body.Password); err != nil {
		return err
	}
	if err := ValidateLoginType(body.Type); err != nil {
		return err
	}
	return nil
}
