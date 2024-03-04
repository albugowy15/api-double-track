package main

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.con/albugowy15/api-double-track/internal/pkg/config"
)

var statement = `
DROP TABLE IF EXISTS schools CASCADE;
DROP TABLE IF EXISTS students CASCADE;
DROP TABLE IF EXISTS admins CASCADE;
DROP TABLE IF EXISTS alternatives CASCADE;
DROP TABLE IF EXISTS questionnare_settings CASCADE;
DROP TABLE IF EXISTS questions CASCADE;
DROP TABLE IF EXISTS answers CASCADE;
DROP TABLE IF EXISTS ahp CASCADE;
DROP TABLE IF EXISTS ahp_to_alternatives CASCADE;
DROP TABLE IF EXISTS ahp_to_alternatives CASCADE;
DROP TABLE IF EXISTS topsis CASCADE;
DROP TABLE IF EXISTS topsis_to_alternatives CASCADE;
DROP FUNCTION IF EXISTS trigger_set_timestamp();
  `

func main() {
	config.LoadConfig(".")
	conf := config.GetConfig()
	connStr := fmt.Sprintf("dbname=%s host=%s port=%s user=%s password=%s sslmode=%s", conf.DbName, conf.DbHost, conf.DbPort, conf.DbUser, conf.DbPass, conf.DbSsl)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	db.MustExec(statement)
	defer db.Close()
}
