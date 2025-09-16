package routes

import (
	"profile-service/api/handlers"
	"profile-service/api/middleware"
	"profile-service/config"
	"profile-service/pkg/user"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func ProfileRouter(app fiber.Router, userservice user.Service, jwtManager *config.JWTManager, v *viper.Viper) {

	profile := app.Group("/user", middleware.GatewayTrust(v))
	profile.Get("", handlers.GetMe(userservice))
	profile.Put("", handlers.UpdateProfile(userservice))

	profile.Post("/link/email", handlers.UpdateEmail(userservice))
	profile.Post("/link/phone", handlers.UpdatePhone(userservice))
}
