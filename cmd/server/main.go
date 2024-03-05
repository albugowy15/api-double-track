package main

import (
	"fmt"
	"log"
	"net/http"

	"github.con/albugowy15/api-double-track/internal/api/router"
	"github.con/albugowy15/api-double-track/internal/pkg/config"
	"github.con/albugowy15/api-double-track/internal/pkg/db"
	"github.con/albugowy15/api-double-track/internal/pkg/utils/jwt"
)

func main() {
	config.LoadConfig(".")
	conf := config.GetConfig()

	db.SetupDB()
	jwt.SetupAuth(conf.Secret)
	api := router.Setup()
	log.Printf("Server running on port %s", conf.Port)
	http.ListenAndServe(":"+conf.Port, api)
	fmt.Println("Hello server")
}
