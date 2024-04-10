package schemas

import (
	"log"
	"math/rand"

	"github.com/albugowy15/api-double-track/internal/pkg/utils"
	"github.com/jmoiron/sqlx"
)

type Admin struct {
	Username    string `db:"username"`
	Password    string `db:"password"`
	Email       string `db:"email"`
	PhoneNumber string `db:"phone_number"`
	SchoolId    string `db:"school_id"`
}

func SeedAdminsTx(tx *sqlx.Tx, schoolIds []SchooldId) {
	// seed admins
	adminPass, err := utils.HashStr("passwordAdmin")
	if err != nil {
		log.Fatalf("error hashing pass: %v", err)
	}
	admins := []Admin{
		{Username: "admintester1", Password: adminPass, Email: "admintester1@gmail.com", PhoneNumber: "087542845123", SchoolId: schoolIds[rand.Intn(3)].Id},
		{Username: "admintester2", Password: adminPass, Email: "admintester2@gmail.com", PhoneNumber: "085166371256", SchoolId: schoolIds[rand.Intn(3)].Id},
		{Username: "admintester3", Password: adminPass, Email: "admintester3@gmail.com", PhoneNumber: "085441327327", SchoolId: schoolIds[rand.Intn(3)].Id},
	}
	_, err = tx.NamedExec(`INSERT INTO admins (username, password, email, phone_number, school_id) VALUES (:username, :password, :email, :phone_number, :school_id)`, admins)
	if err != nil {
		log.Fatalf("error insert admins: %v", err)
	}
}
