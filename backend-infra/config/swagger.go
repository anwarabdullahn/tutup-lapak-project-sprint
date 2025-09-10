package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2"
)

func NewSwagger(app *fiber.App) error {
	// Get the absolute path to ensure it works regardless of working directory
	swaggerPath := "./docs/swagger.json"
	
	// Convert to absolute path but keep relative for flexibility
	absPath, err := filepath.Abs(swaggerPath)
	if err != nil {
		log.Printf("Warning: Could not get absolute path for swagger file: %v", err)
		absPath = swaggerPath // Fallback to relative path
	}
	
	// Check if swagger.json file exists
	if _, err := os.Stat(swaggerPath); os.IsNotExist(err) {
		log.Printf("Warning: Swagger file not found at %s (absolute: %s). Skipping Swagger setup.", swaggerPath, absPath)
		return nil
	}

	log.Printf("Swagger file found at: %s (absolute: %s)", swaggerPath, absPath)

	// Setup Swagger middleware with error recovery
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Error initializing Swagger: %v", r)
		}
	}()

	app.Use(swagger.New(swagger.Config{
		BasePath: "/",
		FilePath: swaggerPath,
		Path:     "swagger",
		Title:    "Tutup Lapak API Documentation",
		CacheAge: 86400,
	}))

	log.Println("✅ Swagger documentation initialized successfully at /swagger")
	return nil
}
