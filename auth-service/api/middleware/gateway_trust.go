package middleware

import (
	"auth-service/api/presenter"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"
)

// GatewayTrust middleware validates requests from the backend-infra gateway
func GatewayTrust(config *viper.Viper) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get the expected internal secret
		internalSecret := config.GetString("INTERNAL_SECRET")
		if internalSecret == "" {
			internalSecret = "backend-infra-internal-secret" // Default for development
		}

		// Check X-Secret header for internal communication security
		xSecret := c.Get("X-Secret")
		if xSecret != internalSecret {
			return c.Status(fiber.StatusUnauthorized).
				JSON(presenter.ErrorResponse("invalid internal secret"))
		}

		// Check X-Auth-Gateway header to ensure it's from our gateway
		xAuthGateway := c.Get("X-Auth-Gateway")
		if xAuthGateway != "backend-infra" {
			return c.Status(fiber.StatusUnauthorized).
				JSON(presenter.ErrorResponse("request not from authorized gateway"))
		}

		// Get user ID from X-User-ID header (set by gateway after JWT validation)
		userID := c.Get("X-User-ID")
		if userID == "" {
			return c.Status(fiber.StatusUnauthorized).
				JSON(presenter.ErrorResponse("missing user context"))
		}

		// Set user context for handlers (same as JWT middleware used to do)
		c.Locals("user_id", userID)
		c.Locals("gateway_validated", true)

		return c.Next()
	}
}
