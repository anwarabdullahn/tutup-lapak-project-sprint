package routes

import (
	"profile-service/api/handlers"
	"profile-service/api/middleware"
	"profile-service/pkg/auth"
	"profile-service/pkg/userfile"

	"github.com/gofiber/fiber/v2"
)

func UserfileRouter(app fiber.Router, userfileService userfile.Service, jm *auth.JWTManager) {

	app.Use(middleware.JWTProtected(jm))

	app.Post(
		"/upload-file",
		handlers.UploadUserFile(userfileService),
	)

}
