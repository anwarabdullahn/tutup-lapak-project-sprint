package routes

import (
	"backend-infra/config"
	"backend-infra/middleware"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// SetupProtectedRoutes sets up JWT-protected routes
func SetupProfileRoutes(app *fiber.App, jwtManager *config.JWTManager) {
	// Protected routes group - all routes here require JWT authentication
	protected := app.Group("/v1/user", middleware.JWTProtected(jwtManager))

	// Auth-related protected routes
	protected.Get("", getUserProfile)

	protected.Put("", UpdateProfile)

	protected.Post("/link/email", UpdateEmail)
	protected.Post("/link/phone", UpdatePhone)
}

// @Summary Get current user profile
// @Description Get authenticated user's profile information
// @Tags         Profile
// @Tags protected
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /v1/user [get]
func getUserProfile(c *fiber.Ctx) error {
	return proxyToProfileService(c, "GET", "/api/v1/user")
}

// UpdatePhone is handler/controller which updates data of current user
// @Summary      Update a phone
// @Description  Update an existing phoe (partial updates allowed)
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        user  body      dtos.PhoneRequest   true  "User update request (partial fields allowed)"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /v1/user/link/phone [post]
func UpdatePhone(c *fiber.Ctx) error {
	return proxyToProfileService(c, "POST", "/api/v1/user/link/phone")
}

// UpdateEmail is handler/controller which updates data of current user
// @Summary      Update a email
// @Description  Update an existing email (partial updates allowed)
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        user  body      dtos.EmailRequest   true  "User update request (partial fields allowed)"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /v1/user/link/email [post]
func UpdateEmail(c *fiber.Ctx) error {
	return proxyToProfileService(c, "POST", "/api/v1/user/link/email")
}

// UpdateProfile is handler/controller which updates data of current user
// @Summary      Update a profile
// @Description  Update an existing profile (partial updates allowed)
// @Tags         Profile
// @Accept       json
// @Produce      json
// @Security BearerAuth
// @Param        user  body      dtos.UpdateUserRequest   true  "User update request (partial fields allowed)"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /v1/user [put]
func UpdateProfile(c *fiber.Ctx) error {
	return proxyToProfileService(c, "PUT", "/api/v1/user")
}

func proxyToProfileService(c *fiber.Ctx, method string, endpoint string) error {
	// Get request body
	body := c.Body()

	// Create request to prfile service
	profile_service_url := c.Locals("service_urls").(*config.ServiceURLs).PROFILE_SERVICE_URL
	url := profile_service_url + endpoint
	fmt.Println("Proxying request to profile service:", url)
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create request",
		})
	}

	// Copy headers
	req.Header.Set("Content-Type", "application/json")
	for key, values := range c.GetReqHeaders() {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}

	// Add user context headers if available (for protected routes)
	if userID := c.Locals("user_id"); userID != nil {
		req.Header.Set("X-User-ID", userID.(string))
		req.Header.Set("X-Auth-Gateway", "backend-infra")

		// Add internal secret for secure communication
		req.Header.Set("X-Secret", "backend-infra-internal-secret") // TODO: Make configurable
	}

	// Make request to auth service
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Auth service unavailable",
		})
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to read response",
		})
	}

	// Set response status and headers
	c.Status(resp.StatusCode)
	for key, values := range resp.Header {
		for _, value := range values {
			c.Set(key, value)
		}
	}

	// Parse and return JSON response
	var jsonResp any
	if err := json.Unmarshal(respBody, &jsonResp); err != nil {
		return c.SendString(string(respBody))
	}

	return c.JSON(jsonResp)
}
