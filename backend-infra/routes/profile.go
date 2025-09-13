package routes

import (
	"backend-infra/config"
	"backend-infra/middleware"
	"log"

	"github.com/gofiber/fiber/v2"
)

// SetupProtectedRoutes sets up JWT-protected routes
func SetupProfileRoutes(app *fiber.App, jwtManager *config.JWTManager) {
	// Protected routes group - all routes here require JWT authentication
	protected := app.Group("/v1/profile", middleware.JWTProtected(jwtManager))

	// Auth-related protected routes
	protected.Get("/me", getUserProfile)
}

// @Summary Get current user profile
// @Description Get authenticated user's profile information
// @Tags protected
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /v1/profile/me [get]
func getUserProfile(c *fiber.Ctx) error {
	log.Println("Fetching user profile")
	return proxyToAuthService(c, "GET", "/api/v1/protected/me")
}
