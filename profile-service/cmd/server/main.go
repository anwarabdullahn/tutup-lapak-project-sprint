package main

import (
	"log"
	"profile-service/config"
)

// @title           Auth Service API
// @version         1.0
// @description     Authentication service for TutupLapak application.

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	v := config.NewViper()
	app := config.NewFiber(v)
	// db := config.NewGorm(v)

	// services := config.InitServices(db)
	// routes.SetupRoutes(app, v, db, services)

	// Run server
	port := v.GetString("PORT")
	if port == "" {
		port = "3001"
	}
	log.Printf("auth-service listening on :%s", port)
	log.Fatal(app.Listen(":" + port))
}
