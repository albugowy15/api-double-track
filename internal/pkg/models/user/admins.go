package user

import (
	"github.com/guregu/null/v5"
)

type Admin struct {
	Id          string      `db:"id"`
	Username    string      `db:"username"`
	Email       null.String `db:"email"`
	Password    string      `db:"password"`
	PhoneNumber null.String `db:"phone_number"`
}
