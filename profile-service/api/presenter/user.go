package presenter

import (
	"profile-service/pkg/entities"

	"github.com/gofiber/fiber/v2"
)

// Response kalau register/login success

// Error response
func ErrorResponse(msg string) string {
	return msg
}

// Success response
func SuccessRegisterResponse(user *entities.User) *fiber.Map {
	return &fiber.Map{
		"ID":    user.ID,
		"Email": user.Email,
	}
}

func ProfileSuccessResponse(data *entities.User) *fiber.Map {
	return &fiber.Map{
		"email": data.Email,
		"phone": data.Phone,
		"fileUri": func() string {
			if data.File != nil {
				return data.File.FileUri
			}
			return ""
		}(),
		"fileThumbnailUri": func() string {
			if data.File != nil {
				return data.File.FileThumbnailUri
			}
			return ""
		}(),
		"bankAccountName":   data.BankAccountName,
		"bankAccountHolder": data.BankAccountHolder,
		"bankAccountNumber": data.BankAccountNumber,
	}
}
