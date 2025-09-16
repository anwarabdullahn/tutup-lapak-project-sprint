package routes

import (
	"backend-infra/config"
	"backend-infra/middleware"
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func SetupFileRoutes(app *fiber.App, jwtManager *config.JWTManager) {
	// Protected routes group - all routes here require JWT authentication
	protected := app.Group("/v1/file", middleware.JWTProtected(jwtManager))

	protected.Post("", UploadFile)
}

// UploadUserFile is handler/controller which upload file
// @Summary      Upload user file
// @Description  Upload user file
// @Tags         Upload File
// @Accept       multipart/form-data
// @Produce      json
// @Security BearerAuth
// @Param        file  formData  file  true  "User File"
// @Success      200   {object}  map[string]interface{}
// @Failure      400   {object}  map[string]interface{}
// @Failure      500   {object}  map[string]interface{}
// @Router       /v1/file [post]
func UploadFile(c *fiber.Ctx) error {
	return proxyToUploadFileService(c, "POST", "/api/v1/file/upload-file")
}

func proxyToUploadFileService(c *fiber.Ctx, method string, endpoint string) error {
	// Get request body
	body := c.Body()

	// Create request to prfile service
	profile_service_url := c.Locals("service_urls").(*config.ServiceURLs).PROFILE_SERVICE_URL
	url := profile_service_url + endpoint
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	// fmt.Println(req)
	// log.Println("Proxying request to:", req.Body)
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
