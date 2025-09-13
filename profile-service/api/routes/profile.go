package routes

import (
	"profile-service/api/handlers"
	"profile-service/api/middleware"
	"profile-service/pkg/auth"
	"profile-service/pkg/user"

	"github.com/gofiber/fiber/v2"
)

func ProfileRouter(app fiber.Router, userservice user.Service, jm *auth.JWTManager) {

	app.Use(middleware.JWTProtected(jm))

	profile := app.Group("/user")
	profile.Get("", handlers.GetMe(userservice))
	profile.Put("", handlers.UpdateProfile(userservice))

	profile.Post("/link/email", handlers.UpdateEmail(userservice))
	profile.Post("/link/phone", handlers.UpdatePhone(userservice))
}
