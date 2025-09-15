// src/api/middleware/jwt.go
package middleware

import (
	"purchase-service/api/presenter"
	"purchase-service/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTProtected(jm *config.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header: "Bearer <token>"
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).
				JSON(presenter.ErrorResponse("missing token"))
		}

		// Delete prefix "Bearer "
		var tokenStr string
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenStr = authHeader[7:]
		} else {
			tokenStr = authHeader
		}

		// Parse token
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			return []byte(jm.Secret()), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).
				JSON(presenter.ErrorResponse("invalid or expired token"))
		}

		// Extract claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Locals("user_id", claims["user_id"])
		}

		return c.Next()
	}
}
