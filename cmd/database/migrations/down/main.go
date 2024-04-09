package main

import (
	"github.com/albugowy15/api-double-track/internal/pkg/config"
	"github.com/albugowy15/api-double-track/internal/pkg/db"
	_ "github.com/lib/pq"
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
	db.GetDb().MustExec(statement)
	db.GetDb().Close()
}
