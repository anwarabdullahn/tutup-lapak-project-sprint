package main

import (
	"log"

	"purchase-service/api/routes"
	"purchase-service/config"
)

// @title           Purchase Service API
// @version         1.0
// @description     Purchase service for TutupLapak application.

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	v := config.NewViper()
	app := config.NewFiber(v)
	db := config.NewGorm(v)

	services := config.InitServices(db)
	routes.SetupRoutes(app, v, db, services)

	// Run server
	port := v.GetString("SERVER_PORT")
	if port == "" {
		port = "3004"
	}
	log.Printf("purchase-service listening on :%s", port)
	log.Fatal(app.Listen(":" + port))
}
