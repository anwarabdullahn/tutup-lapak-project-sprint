package presenter

import (
	"auth-service/pkg/entities"

	"github.com/gofiber/fiber/v2"
)

// Response kalau register/login success
func UserSuccessResponse(user *entities.User) *fiber.Map {
	return &fiber.Map{
		"email": user.Email,
	}
}

// Error response
func ErrorResponse(msg string) string {
	return msg
}
