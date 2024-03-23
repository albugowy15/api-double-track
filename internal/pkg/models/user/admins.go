package user

import (
	"github.com/guregu/null/v5"
)

type Admin struct {
	Id          string      `db:"id" json:"id"`
	Username    string      `db:"username" json:"username"`
	Email       null.String `db:"email" json:"email,omitempty"`
	Password    string      `db:"password" json:"-"`
	PhoneNumber null.String `db:"phone_number" json:"phone_number,omitempty"`
}

type UpdateAdminRequest struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}
