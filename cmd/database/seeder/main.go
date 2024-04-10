package main

import (
	"log"

	"github.com/albugowy15/api-double-track/cmd/database/seeder/schemas"
	"github.com/albugowy15/api-double-track/internal/pkg/config"
	"github.com/albugowy15/api-double-track/internal/pkg/db"
	_ "github.com/lib/pq"
)

func init() {
	config.LoadConfig(".")
	db.SetupDB()
}

func main() {
	env := config.GetConfig().AppEnv
	tx := db.GetDb().MustBegin()

	schemas.SeedAlternativeTx(tx)
	schemas.SeedQuestionsTx(tx)

	if env == "prod" {
		log.Println("seed database for production")
		tx.Commit()
		db.GetDb().Close()
		return
	}

	log.Println("seed database for development")
	schemas.SeedSchools()
	schoolIds := schemas.SeedStudentsTx(tx)
	schemas.SeedAdminsTx(tx, schoolIds)

	tx.Commit()
	db.GetDb().Close()
}
