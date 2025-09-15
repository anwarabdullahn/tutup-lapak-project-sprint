package routes

import (
	"profile-service/api/handlers"
	"profile-service/api/middleware"
	"profile-service/config"
	"profile-service/pkg/uploadfile"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

func UploadfileRouter(app fiber.Router, userService uploadfile.Service, jwtManager *config.JWTManager, v *viper.Viper) {

	file := app.Group("/file", middleware.GatewayTrust(v))

	file.Post(
		"/upload-file",
		handlers.UploadFile(userService),
	)

}
