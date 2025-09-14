package routes

import (
	"purchase-service/config"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

// SetupRoutes configures all application routes
func SetupRoutes(app *fiber.App, v *viper.Viper, db *gorm.DB, services config.Services) {
	// API v1 group
	api := app.Group("/api/v1")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "purchase-service",
			"status":  "running",
		})
	})

	PurchaseRouter(api, services)

	app.Get("/healthz", func(c *fiber.Ctx) error {
		sqlDB, err := db.DB() // get underlying *sql.DB from GORM
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Database connection error")
		}

		if err := sqlDB.Ping(); err != nil { // try pinging the DB
			return c.Status(fiber.StatusInternalServerError).SendString("Database not reachable")
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "ok",
		}) // 200 OK if DB is fine
	})
}

