package validator

import (
	"errors"
	"net/mail"
	"strconv"
	"strings"

	"github.com/albugowy15/api-double-track/internal/pkg/models/user"
)

var (
	ErrFullnameEmpty        = errors.New("nama lengkap wajib diisi")
	ErrEmailEmpty           = errors.New("email wajib diisi")
	ErrEmailInvalid         = errors.New("email tidak valid")
	ErrNisnEmpty            = errors.New("nisn wajib diisi")
	ErrNisnInvalid          = errors.New("nisn wajib berupa angka")
	ErrPhoneNumberLength    = errors.New("nomor hp hanya boleh terdiri dari 10 sampai 14 digit angka")
	ErrPhoneNumberPrefix    = errors.New("nomor hp diawali dengan 08")
	ErrPhoneNumberNotNumber = errors.New("nomor hp hanya boleh terdiri dari angka")
)

func ValidateFullname(fullname string) error {
	if len(fullname) == 0 {
		return ErrFullnameEmpty
	}
	return nil
}

func ValidateEmail(email string) error {
	if len(email) == 0 {
		return ErrEmailEmpty
	}
	_, err := mail.ParseAddress(email)
	if err != nil {
		return ErrEmailInvalid
	}
	return nil
}

func ValidateNisn(nisn string) error {
	if len(nisn) == 0 {
		return ErrNisnEmpty
	}
	_, err := strconv.Atoi(nisn)
	if err != nil {
		return ErrNisnInvalid
	}
	return nil
}

func ValidatePhoneNumber(phoneNumber string) error {
	if len(phoneNumber) < 10 || len(phoneNumber) > 14 {
		return ErrPhoneNumberLength
	}
	if !strings.HasPrefix(phoneNumber, "08") {
		return ErrPhoneNumberPrefix
	}
	_, err := strconv.Atoi(phoneNumber)
	if err != nil {
		return ErrPhoneNumberNotNumber
	}
	return nil
}

func ValidateAddStudent(data user.Student) (user.Student, error) {
	var sanitized user.Student
	if err := ValidateFullname(data.Fullname); err != nil {
		return sanitized, err
	}
	if err := ValidateNisn(data.Nisn); err != nil {
		return sanitized, err
	}
	nisn := strings.TrimSpace(data.Nisn)
	fullname := strings.TrimSpace(data.Fullname)
	sanitized.Nisn = nisn
	sanitized.Fullname = fullname

	return sanitized, nil
}

func ValidateUpdateStudent(data user.Student) (user.Student, error) {
	var sanitized user.Student

	sanitized.Fullname = strings.TrimSpace(data.Fullname)
	sanitized.Username = strings.TrimSpace(data.Username)
	sanitized.Email.Valid = data.Email.Valid
	sanitized.Email.String = strings.TrimSpace(data.Email.String)
	sanitized.Nisn = strings.TrimSpace(data.Nisn)
	sanitized.PhoneNumber.Valid = data.PhoneNumber.Valid
	sanitized.PhoneNumber.String = strings.TrimSpace(data.PhoneNumber.String)

	if err := ValidateFullname(sanitized.Fullname); err != nil {
		return sanitized, err
	}
	if err := ValidateUsername(sanitized.Username); err != nil {
		return sanitized, err
	}
	if sanitized.Email.Valid {
		if err := ValidateEmail(sanitized.Email.ValueOrZero()); err != nil {
			return sanitized, err
		}
	}
	if err := ValidateNisn(sanitized.Nisn); err != nil {
		return sanitized, err
	}
	if sanitized.PhoneNumber.Valid {
		if err := ValidatePhoneNumber(sanitized.PhoneNumber.ValueOrZero()); err != nil {
			return sanitized, err
		}
	}
	return sanitized, nil
}
