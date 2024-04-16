package db

import (
	"log"

	"github.com/albugowy15/api-double-track/pkg/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var AppDB *sqlx.DB

func Init() {
	config.Load(".")

	db, err := sqlx.Connect("postgres", config.AppConfig.DatabaseUrl)
	if err != nil {
		log.Fatalf("error connecting with database: %v", err)
	}
	AppDB = db
}
