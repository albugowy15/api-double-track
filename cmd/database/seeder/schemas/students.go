package schemas

import (
	"log"
	"math/rand"

	"github.com/albugowy15/api-double-track/internal/pkg/db"
	"github.com/albugowy15/api-double-track/internal/pkg/utils"
	"github.com/jmoiron/sqlx"
)

type Student struct {
	Username    string `db:"username"`
	Password    string `db:"password"`
	Fullname    string `db:"fullname"`
	Nisn        string `db:"nisn"`
	Email       string `db:"email"`
	PhoneNumber string `db:"phone_number"`
	SchoolId    string `db:"school_id"`
}

type SchooldId struct {
	Id string `db:"id"`
}

func SeedStudentsTx(tx *sqlx.Tx) []SchooldId {
	schoolIds := []SchooldId{}
	db.GetDb().Select(&schoolIds, "SELECT id from schools LIMIT 3")
	// students
	studentPass, err := utils.HashStr("passwordStudent")
	if err != nil {
		log.Fatalf("error hashing pass: %v", err)
		return schoolIds
	}
	students := []Student{
		{
			Username:    "bughowy",
			Email:       "bughowy@gmail.com",
			Password:    studentPass,
			Fullname:    "Mohamad Kholid Bughowi",
			Nisn:        "123252532",
			PhoneNumber: "086345193034",
			SchoolId:    schoolIds[rand.Intn(3)].Id,
		},
		{
			Username:    "ahmadthoriq",
			Email:       "ahmadthoriq@gmail.com",
			Password:    studentPass,
			Fullname:    "Ahmad Thoriq",
			Nisn:        "12773663943",
			PhoneNumber: "084495723723",
			SchoolId:    schoolIds[rand.Intn(3)].Id,
		},
		{
			Username:    "dwiassegaf",
			Email:       "dwiassegaf@gmail.com",
			Password:    studentPass,
			Fullname:    "Dwi Assegaf",
			Nisn:        "1344263626",
			PhoneNumber: "076653626642",
			SchoolId:    schoolIds[rand.Intn(3)].Id,
		},
		{
			Username:    "efendi",
			Email:       "efendimanik@gmail.com",
			Password:    studentPass,
			Fullname:    "Efendi Manik",
			Nisn:        "1283774334",
			PhoneNumber: "08423173732",
			SchoolId:    schoolIds[rand.Intn(3)].Id,
		},
		{
			Username:    "gedesam",
			Email:       "gedesam@gmail.com",
			Password:    studentPass,
			Fullname:    "Gede Samudra",
			Nisn:        "1264378882",
			PhoneNumber: "0816763434999",
			SchoolId:    schoolIds[rand.Intn(3)].Id,
		},
		{
			Username:    "gatotbakti",
			Email:       "gatotbakti@gmail.com",
			Password:    studentPass,
			Fullname:    "Gatot Surbakti",
			Nisn:        "74346364343",
			PhoneNumber: "0857285926385",
			SchoolId:    schoolIds[rand.Intn(3)].Id,
		},
		{
			Username:    "harissag",
			Email:       "harissag@gmail.com",
			Password:    studentPass,
			Fullname:    "Haris Saragih",
			Nisn:        "12665577557",
			PhoneNumber: "0865523558834",
			SchoolId:    schoolIds[rand.Intn(3)].Id,
		},
	}
	_, err = tx.NamedExec("INSERT INTO students (username, email, password, fullname, nisn, phone_number, school_id) VALUES (:username, :email, :password, :fullname, :nisn, :phone_number, :school_id)", students)
	if err != nil {
		log.Fatalf("error insert students: %v", err)
		return schoolIds
	}
	return schoolIds
}
