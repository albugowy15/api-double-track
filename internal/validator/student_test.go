package validator_test

import (
	"errors"
	"testing"

	"github.com/albugowy15/api-double-track/internal/models"
	"github.com/albugowy15/api-double-track/internal/validator"
)

func TestValidateAddStudent(t *testing.T) {
	body := models.Student{
		Fullname: "",
		Nisn:     "",
	}
	_, err := validator.ValidateAddStudent(body)
	if err != validator.ErrFullnameEmpty {
		t.Errorf("error not match got: %v", err)
	}

	body.Fullname = "Ahmad dhhahjjjaf"
	_, err = validator.ValidateAddStudent(body)
	if err != validator.ErrNisnEmpty {
		t.Errorf("error not match got: %v", err)
	}

	body.Nisn = "12512515247743"
	body.Fullname = "  Ahmed    "
	sanit, err := validator.ValidateAddStudent(body)
	if err != nil {
		t.Errorf("expect no error, got: %v", err)
	}
	if sanit.Fullname != "Ahmed" {
		t.Errorf("expect fullname santized, got: %s", sanit.Fullname)
	}

	body.Nisn = "236362H   H&&2352352"
	sanit, err = validator.ValidateAddStudent(body)
	if err != validator.ErrNisnInvalid {
		t.Errorf("error not match got: %v", err)
	}

	body.Nisn = "123344545454"
	sanit, err = validator.ValidateAddStudent(body)
	if err != nil {
		t.Errorf("expect no error, got %v", err)
	}
}

func TestValidateUpdateStudent(t *testing.T) {
	data := models.Student{
		Fullname: "Mohamad kholid",
		Nisn:     "14232332",
		Username: "feuufe77232jj",
	}
	_, err := validator.ValidateUpdateStudent(data)
	if err != nil {
		t.Errorf("expect no error, got: %v", err)
	}

	data.Fullname = ""
	_, err = validator.ValidateUpdateStudent(data)
	if err != validator.ErrFullnameEmpty {
		t.Errorf("expect error: %v, got: %v", errors.New("nama lengkap wajib diisi"), err)
	}

	data.Fullname = "Mohamad Kholid"

	data.Nisn = "2232 fefe 773743"
	_, err = validator.ValidateUpdateStudent(data)
	if err != validator.ErrNisnInvalid {
		t.Errorf("expect error: %v, got: %v", errors.New("nisn wajib berupa angka"), err)
	}

	data.Nisn = ""
	_, err = validator.ValidateUpdateStudent(data)
	if err != validator.ErrNisnEmpty {
		t.Errorf("expect error: %v, got: %v", errors.New("nisn wajib diisi"), err)
	}
	data.Nisn = "12164664343"

	data.Username = ""
	_, err = validator.ValidateUpdateStudent(data)
	if err.Error() != "username harus terdiri dari 6 hingga 30 karakter" {
		t.Errorf("expect error: %v, got: %v", errors.New("username harus terdiri dari 6 hingga 30 karakter"), err)
	}
	data.Username = "realUser123"

	data.Email.Valid = true
	data.Email.String = "kholdi"
	_, err = validator.ValidateUpdateStudent(data)
	if err != validator.ErrEmailInvalid {
		t.Errorf("expect error: %v, got: %v", errors.New("email tidak valid"), err)
	}
	data.Email.String = "kholidbugh@gmail.com"

	data.PhoneNumber.Valid = true
	data.PhoneNumber.String = "7434636"
	_, err = validator.ValidateUpdateStudent(data)
	if err != validator.ErrPhoneNumberLength {
		t.Errorf("expect error: %v, got: %v", errors.New("nomor hp hanya boleh terdiri dari 10 sampai 14 digit angka"), err)
	}

	data.PhoneNumber.String = "12345678910"
	_, err = validator.ValidateUpdateStudent(data)
	if err != validator.ErrPhoneNumberPrefix {
		t.Errorf("expect error: %v, got: %v", errors.New("nomor hp diawali dengan 08"), err)
	}

	data.PhoneNumber.String = "08123456789we"
	_, err = validator.ValidateUpdateStudent(data)
	if err != validator.ErrPhoneNumberNotNumber {
		t.Errorf("expect error: %v, got: %v", errors.New("nomor hp hanya boleh terdiri dari angka"), err)
	}

	data.Fullname = "     Mohamaed Kholid    "
	data.Username = "gfetf6"
	data.Email.String = "    kholidnfefe@gmail.com   "
	data.PhoneNumber.String = "    08883232323   "
	sanit, err := validator.ValidateUpdateStudent(data)
	if err != nil {
		t.Errorf("expect no error, got: %v", err)
	}
	if sanit.Fullname != "Mohamaed Kholid" && sanit.Username != "gfetf6" && sanit.Email.String != "kholidnfefe@gmail.com" && sanit.PhoneNumber.String != "08883232323" {
		t.Errorf("expect data got sanitized, got none")
	}
}
