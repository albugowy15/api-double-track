package main

import (
	"log"
	"net/http"

	"github.com/albugowy15/api-double-track/db"
	"github.com/albugowy15/api-double-track/docs"
	"github.com/albugowy15/api-double-track/internal/router"
	"github.com/albugowy15/api-double-track/internal/services"
	"github.com/albugowy15/api-double-track/pkg/auth"
	"github.com/albugowy15/api-double-track/pkg/config"
)

func init() {
	config.Load(".")
	db.Init()
	services.InitSubCriteriaWeights()
}

//	@title			Double Track API
//	@version		1.0
//	@description	This is a Double Track REST API

// @BasePath	/v1
func main() {
	auth.Init(config.AppConfig.Secret)
	api := router.Init()
	docs.SwaggerInfo.Host = config.AppConfig.BaseUrl
	log.Printf("Server running on port %s", config.AppConfig.Port)
	log.Printf("See the docs: %s/swagger/index.html", docs.SwaggerInfo.Host)
	http.ListenAndServe(":"+config.AppConfig.Port, api)
}
