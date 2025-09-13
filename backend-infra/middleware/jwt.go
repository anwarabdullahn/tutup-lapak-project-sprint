package middleware

import (
	"backend-infra/config"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWTProtected middleware validates JWT tokens at the gateway level
func JWTProtected(jwtManager *config.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header: "Bearer <token>"
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "missing authorization token",
				"error":   "authorization header is required",
			})
		}

		// Extract token from "Bearer <token>" format
		var tokenStr string
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenStr = authHeader[7:]
		} else {
			tokenStr = authHeader
		}

		// Parse and validate token
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			// Ensure the signing method is HMAC
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "invalid signing method")
			}
			return []byte(jwtManager.Secret()), nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "invalid token",
				"error":   err.Error(),
			})
		}

		if !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "token is not valid",
				"error":   "token validation failed",
			})
		}

		// Extract user information from token claims
		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			userID, exists := claims["user_id"]
			if !exists {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"success": false,
					"message": "invalid token claims",
					"error":   "user_id not found in token",
				})
			}

			// Store user context for use in handlers
			c.Locals("user_id", userID)
			c.Locals("jwt_claims", claims)

			// Add user context headers for downstream services
			c.Set("X-User-ID", userID.(string))
			c.Set("X-Auth-Gateway", "backend-infra")
		}

		return c.Next()
	}
}