package db

import (
	"log"

	"github.com/albugowy15/api-double-track/internal/pkg/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var dbConn *sqlx.DB

func SetupDB() {
	config.LoadConfig(".")
	conf := config.GetConfig()

	db, err := sqlx.Connect("postgres", conf.DatabaseUrl)
	if err != nil {
		log.Fatalf("error connecting with database: %v", err)
	}
	dbConn = db
}

func GetDb() *sqlx.DB {
	return dbConn
}
