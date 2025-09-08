package routes

import (
	"auth-service/api/handlers"
	"auth-service/api/middleware"
	"auth-service/pkg/user"

	"github.com/gofiber/fiber/v2"
)

func UserRouter(app fiber.Router, service user.Service, jwtManager *user.JWTManager) {
	app.Post("/register", handlers.Register(service))
	app.Post("/login", handlers.Login(service, jwtManager))

	protected := app.Group("/protected", middleware.JWTProtected(jwtManager))
	protected.Get("/me", func(c *fiber.Ctx) error {
		userID := c.Locals("user_id")
		return c.JSON(fiber.Map{
			"message": "success",
			"user_id": userID,
		})
	})
}
