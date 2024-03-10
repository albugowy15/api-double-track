package db

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.con/albugowy15/api-double-track/internal/pkg/config"
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
