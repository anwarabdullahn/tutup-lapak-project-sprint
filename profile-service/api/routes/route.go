package routes

import (
	"profile-service/config"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, v *viper.Viper, db *gorm.DB, services config.Services) {
	// API v1 group
	api := app.Group("/api/v1")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": "auth-service",
			"status":  "running",
		})
	})

	// Init JWT Manager (24 jam expired)
	jwtManager := config.NewJWTManager(v.GetString("JWT_SECRET"), 24*time.Hour)

	ProfileRouter(api, services.UserService, jwtManager, v)
	UploadfileRouter(api, services.FileService, jwtManager, v)

}
