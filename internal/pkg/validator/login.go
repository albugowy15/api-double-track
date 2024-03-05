package validator

import (
	"fmt"
	"regexp"
	"strings"

	"github.con/albugowy15/api-double-track/internal/pkg/models"
)

func ValidateUsername(username string) error {
	usernameLen := len(username)
	if usernameLen < 6 || usernameLen > 30 {
		return fmt.Errorf("username harus terdiri dari 6 hingga 30 karakter")
	}
	pattern := "^[a-zA-Z0-9]+$"
	re := regexp.MustCompile(pattern)
	if !re.MatchString(username) {
		return fmt.Errorf("username hanya boleh terdiri dari huruf alpabet atau angka tanpa spasi")
	}
	return nil
}

func ValidatePassword(password string) error {
	passwordLength := len(password)
	if passwordLength < 8 || passwordLength > 16 {
		return fmt.Errorf("password harus terdiri dari 8 hingga 16 karakter")
	}
	if strings.Contains(password, " ") {
		return fmt.Errorf("password tidak boleh terdapat spasi")
	}
	upperRegex := `[A-Z]`
	if !regexp.MustCompile(upperRegex).MatchString(password) {
		return fmt.Errorf("password minimal terdapat 1 huruf kapital")
	}
	lowerRegex := `[a-z]`
	if !regexp.MustCompile(lowerRegex).MatchString(password) {
		return fmt.Errorf("password minimal terdapat 1 huruf kecil")
	}
	// digitRegex := `[0-9]`
	// if !regexp.MustCompile(digitRegex).MatchString(password) {
	// 	return fmt.Errorf("password minimal terdapat 1 angka")
	// }
	return nil
}

func ValidateLoginType(loginType string) error {
	if loginType != "admin" && loginType != "student" {
		return fmt.Errorf("type login tidak valid")
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
