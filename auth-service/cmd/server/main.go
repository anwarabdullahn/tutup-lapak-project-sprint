package main

import (
	"log"

	"auth-service/config"
	"github.com/gofiber/fiber/v2"
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

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "auth-service",
			"status":  "running",
		})
	})

	// Run server
	port := v.GetString("PORT")
	if port == "" {
		port = "3001"
	}
	log.Printf("auth-service listening on :%s", port)
	log.Fatal(app.Listen(":" + port))
}
