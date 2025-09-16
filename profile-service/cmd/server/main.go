package main

import (
	"log"
	"profile-service/api/routes"
	"profile-service/config"

	"github.com/gofiber/contrib/swagger"
)

// @title My API
// @version 1.0
// @description This is my API with JWT auth

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {

	v := config.NewViper()
	app := config.NewFiber(v)
	db := config.NewGorm(v)

	services := config.InitServices(db)
	routes.SetupRoutes(app, v, db, services)

	app.Use(swagger.New(swagger.Config{
		BasePath: "/",
		FilePath: "./docs/swagger.json",
		Path:     "swagger",
		Title:    "Swagger API Docs",
		CacheAge: 86400,
	}))

	// Run server
	port := v.GetString("server.port")
	if port == "" {
		port = "3002"
	}
	log.Fatal(app.Listen(":" + port))
}
