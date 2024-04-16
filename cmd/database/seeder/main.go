package main

import (
	"log"

	"github.com/albugowy15/api-double-track/db"
	"github.com/albugowy15/api-double-track/db/seeds"
	"github.com/albugowy15/api-double-track/pkg/config"
	_ "github.com/lib/pq"
)

func init() {
	config.Load(".")
	db.Init()
}

func main() {
	env := config.AppConfig.AppEnv
	tx := db.AppDB.MustBegin()

	seeds.SeedAlternativeTx(tx)
	seeds.SeedQuestionsTx(tx)

	if env == "prod" {
		log.Println("seed database for production")
		tx.Commit()
		db.AppDB.Close()
		return
	}

	log.Println("seed database for development")
	seeds.SeedSchools()
	schoolIds := seeds.SeedStudentsTx(tx)
	seeds.SeedAdminsTx(tx, schoolIds)

	tx.Commit()
	db.AppDB.Close()
}
