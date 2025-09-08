package routes

import (
	"auth-service/api/handlers"
	"auth-service/api/middleware"
	"auth-service/config"
	"auth-service/pkg/user"

	"github.com/gofiber/fiber/v2"
)

func UserRouter(app fiber.Router, service user.Service, jwtManager *config.JWTManager) {
	// Registration endpoints
	app.Post("/register/email", handlers.RegisterEmail(service))
	app.Post("/register/phone", handlers.RegisterPhone(service))

	// Login endpoints
	app.Post("/login/email", handlers.LoginEmail(service, jwtManager))
	app.Post("/login/phone", handlers.LoginPhone(service, jwtManager))

	// Protected routes
	protected := app.Group("/protected", middleware.JWTProtected(jwtManager))
	protected.Get("/me", func(c *fiber.Ctx) error {
		userID := c.Locals("user_id")
		return c.JSON(fiber.Map{
			"message": "success",
			"user_id": userID,
		})
	})
}
