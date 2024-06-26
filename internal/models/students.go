package models

import "github.com/guregu/null/v5"

type Student struct {
	Id          string      `db:"id" json:"id"`
	Username    string      `db:"username" json:"username"`
	Password    string      `db:"password" json:"-"`
	Fullname    string      `db:"fullname" json:"fullname"`
	Nisn        string      `db:"nisn" json:"nisn"`
	Email       null.String `db:"email" json:"email,omitempty"`
	PhoneNumber null.String `db:"phone_number" json:"phone_number,omitempty"`
}

type StudentProfile struct {
	Id          string      `db:"id" json:"id"`
	Username    string      `db:"username" json:"username"`
	Fullname    string      `db:"fullname" json:"fullname"`
	Nisn        string      `db:"nisn" json:"nisn"`
	School      string      `json:"school"`
	Email       null.String `db:"email" json:"email,omitempty"`
	PhoneNumber null.String `db:"phone_number" json:"phone_number,omitempty"`
}

type StudentRegisterRequest struct {
	Fullname string `json:"fullname"`
	Nisn     string `json:"nisn"`
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
	School   string `json:"school"`
}
