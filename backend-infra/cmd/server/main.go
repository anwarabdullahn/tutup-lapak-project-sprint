package main

import (
	"log"

	"backend-infra/config"
	"backend-infra/routes"

	"github.com/gofiber/fiber/v2"
)

// @title           TutupLapak API
// @version         1.0
// @description     This is a sample server for a tutup lapak application.

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	v := config.NewViper()
	app := config.NewFiber(v)

	if err := config.NewSwagger(app); err != nil {
		log.Printf("Failed to initialize Swagger: %v", err)
	}

	app.Get("/healthz", func(c *fiber.Ctx) error { return c.SendString("ok") })

	// Setup routes
	routes.SetupAuthRoutes(app)

	// Run server
	port := v.GetString("SERVER_PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}
