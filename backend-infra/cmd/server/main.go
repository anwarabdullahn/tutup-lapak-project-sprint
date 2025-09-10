package main

import (
	"log"
	"time"

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

	// Initialize JWT Manager
	jwtSecret := v.GetString("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "nv6FNtvAmBmUMHRSta8aSZNwiw4XAH" // Same as auth-service default
		log.Println("Warning: Using default JWT secret. Set JWT_SECRET environment variable.")
	}
	jwtManager := config.NewJWTManager(jwtSecret, 24*time.Hour)

	// Setup routes
	routes.SetupAuthRoutes(app)
	routes.SetupProfileRoutes(app, jwtManager)

	// Run server
	port := v.GetString("SERVER_PORT")
	if port == "" {
		port = "3000"
	}
	log.Fatal(app.Listen(":" + port))
}
