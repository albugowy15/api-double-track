package db

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.con/albugowy15/api-double-track/internal/pkg/config"
)

var dbConn *sqlx.DB

func SetupDB() {
	config.LoadConfig(".")
	conf := config.GetConfig()
	connStr := fmt.Sprintf("dbname=%s host=%s port=%s user=%s password=%s sslmode=%s", conf.DbName, conf.DbHost, conf.DbPort, conf.DbUser, conf.DbPass, conf.DbSsl)

	db, err := sqlx.Connect(conf.DbDriver, connStr)
	if err != nil {
		log.Fatalf("error connecting with database: %v", err)
	}
	dbConn = db
}

func GetDb() *sqlx.DB {
	return dbConn
}
