package main

import (
	"log"

	_ "github.com/lib/pq"
	"github.con/albugowy15/api-double-track/internal/pkg/config"
	"github.con/albugowy15/api-double-track/internal/pkg/db"
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

func init() {
	config.LoadConfig(".")
	db.SetupDB()
}

func main() {
	if config.GetConfig().AppEnv == "prod" {
		log.Fatalf("you cannot run database migrations when app is running in production")
	}
	db.GetDb().MustExec(statement)
	db.GetDb().Close()
}
