package main

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func main() {
	db, err := sqlx.Connect("postgres", "user=postgres password=postgres dbname=dbtrack sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
}
