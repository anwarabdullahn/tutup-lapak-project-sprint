package routes

import (
	"auth-service/api/handlers"
	"auth-service/api/middleware"
	"auth-service/config"
	"auth-service/pkg/user"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func UserRouter(app fiber.Router, service user.Service, jwtManager *config.JWTManager, v *viper.Viper) {
	// Registration endpoints
	app.Post("/register/email", handlers.RegisterEmail(service))
	app.Post("/register/phone", handlers.RegisterPhone(service))

	// Login endpoints
	app.Post("/login/email", handlers.LoginEmail(service, jwtManager))
	app.Post("/login/phone", handlers.LoginPhone(service, jwtManager))

	// Protected routes - now trust the gateway instead of validating JWT directly
	protected := app.Group("/protected", middleware.GatewayTrust(v))

	// Add protected endpoints
	protected.Get("/me", handlers.GetMe(service))
}
