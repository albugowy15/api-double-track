package seeds

import (
	"log"

	"github.com/albugowy15/api-double-track/db"
)

type School struct {
	Name string `db:"name"`
}

func SeedSchools() {
	// seed schools
	schools := []School{
		{Name: "SMA IPIEMS Surabaya"},
		{Name: "SMA Dharmawanita Surabaya"},
		{Name: "SMA Negeri 1 Ngadirojo"},
		{Name: "SMA Negeri 1 Jenangan"},
		{Name: "SMA Negeri 1 Gondanglegi"},
		{Name: "SMA Negeri 1 Balongpanggang"},
		{Name: "SMA Negeri 1 Turen"},
		{Name: "SMA Negeri 1 Sumbermanjing"},
		{Name: "SMA Negeri 1 Pulung"},
	}
	_, err := db.AppDB.NamedExec(`INSERT INTO schools (name) VALUES (:name)`, schools)
	if err != nil {
		log.Fatalf("error insert schools: %v", err)
	}
}
