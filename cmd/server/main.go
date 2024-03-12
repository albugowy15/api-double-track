package main

import (
	"log"
	"net/http"

	"github.com/albugowy15/api-double-track/internal/api/router"
	"github.com/albugowy15/api-double-track/internal/pkg/config"
	"github.com/albugowy15/api-double-track/internal/pkg/db"
	"github.com/albugowy15/api-double-track/internal/pkg/utils/jwt"
)

func init() {
	config.LoadConfig(".")
	db.SetupDB()
}

func main() {
	conf := config.GetConfig()
	jwt.SetupAuth(conf.Secret)
	api := router.Setup()
	log.Printf("Server running on port %s", conf.Port)
	http.ListenAndServe(":"+conf.Port, api)
}
