package routes

import (
	"backend-infra/config"
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

// SetupAuthRoutes sets up auth-related routes
func SetupAuthRoutes(app *fiber.App) {
	v1 := app.Group("/v1")

	// Login routes
	v1.Post("/login/email", loginWithEmail)
	v1.Post("/login/phone", loginWithPhone)

	// Register routes
	v1.Post("/register/email", registerWithEmail)
	v1.Post("/register/phone", registerWithPhone)
}

// @Summary Login with email
// @Description Login user with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dtos.EmailLoginRequest true "Login credentials"
// @Success 200 {object} dtos.LoginResponse
// @Failure 400 {object} map[string]interface{}
// @Router /v1/login/email [post]
func loginWithEmail(c *fiber.Ctx) error {
	return proxyToAuthService(c, "POST", "/api/v1/login/email")
}

// @Summary Login with phone
// @Description Login user with phone and password
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dtos.PhoneLoginRequest true "Login credentials"
// @Success 200 {object} dtos.LoginResponse
// @Failure 400 {object} map[string]interface{}
// @Router /v1/login/phone [post]
func loginWithPhone(c *fiber.Ctx) error {
	return proxyToAuthService(c, "POST", "/api/v1/login/phone")
}

// @Summary Register with email
// @Description Register new user with email
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dtos.EmailRegisterRequest true "Registration data"
// @Success 201 {object} dtos.LoginResponse
// @Failure 400 {object} map[string]interface{}
// @Router /v1/register/email [post]
func registerWithEmail(c *fiber.Ctx) error {
	return proxyToAuthService(c, "POST", "/api/v1/register/email")
}

// @Summary Register with phone
// @Description Register new user with phone
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dtos.PhoneRegisterRequest true "Registration data"
// @Success 201 {object} dtos.LoginResponse
// @Failure 400 {object} map[string]interface{}
// @Router /v1/register/phone [post]
func registerWithPhone(c *fiber.Ctx) error {
	return proxyToAuthService(c, "POST", "/api/v1/register/phone")
}

// proxyToAuthService forwards requests to the auth service
func proxyToAuthService(c *fiber.Ctx, method string, endpoint string) error {
	// Get request body
	body := c.Body()

	// Create request to auth service
	auth_service_url := c.Locals("service_urls").(*config.ServiceURLs).AuthServiceURL
	url := auth_service_url + endpoint
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
