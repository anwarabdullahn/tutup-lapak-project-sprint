package routes

import (
	"profile-service/api/handlers"
	"profile-service/pkg/auth"

	"github.com/gofiber/fiber/v2"
)

func AuthRouter(app fiber.Router, service auth.Service, jwtManager *auth.JWTManager) {

	app.Post("/register", handlers.Register(service))
	app.Post("/login", handlers.Login(service, jwtManager))

	// protected := app.Group("/protected", middleware.JWTProtected(jwtManager))
	// protected.Get("/me", func(c *fiber.Ctx) error {
	// 	userID := c.Locals("user_id")
	// 	user := service.FindByID(userID.(string))
	// 	return c.JSON(fiber.Map{
	// 		"message": "success",
	// 		"user_id": userID,
	// 	})
	// })
}
