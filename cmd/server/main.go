package main

import (
	"log"
	"net/http"

	"github.com/albugowy15/api-double-track/docs"
	"github.com/albugowy15/api-double-track/internal/api/router"
	"github.com/albugowy15/api-double-track/internal/api/services"
	"github.com/albugowy15/api-double-track/internal/pkg/config"
	"github.com/albugowy15/api-double-track/internal/pkg/db"
	"github.com/albugowy15/api-double-track/internal/pkg/utils/jwt"
)

func init() {
	config.LoadConfig(".")
	db.SetupDB()
	services.InitSubCriteriaWeights()
}

//	@title			Double Track API
//	@version		1.0
//	@description	This is a Double Track REST API

// @BasePath	/v1
func main() {
	conf := config.GetConfig()
	jwt.SetupAuth(conf.Secret)
	api := router.Setup()
	docs.SwaggerInfo.Host = conf.BaseUrl
	log.Printf("Server running on port %s", conf.Port)
	log.Printf("See the docs: %s/swagger/index.html", docs.SwaggerInfo.Host)
	http.ListenAndServe(":"+conf.Port, api)
}
