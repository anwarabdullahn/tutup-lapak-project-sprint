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

// SetupPurchaseRoutes sets up purchase-related routes
func SetupPurchaseRoutes(app *fiber.App, jwtManager *config.JWTManager) {
	// Protected routes group - all routes here require JWT authentication
	protected := app.Group("/v1/purchase", middleware.JWTProtected(jwtManager))

	// Purchase routes
	protected.Post("/", createPurchase)
	protected.Get("/", listPurchases)
	protected.Get("/:purchaseId", getPurchaseByID)
	protected.Post("/:purchaseId", uploadPaymentProof)
}

// @Summary Create a new purchase
// @Description Customer can add their items to cart so they can pay them
// @Tags purchase
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dtos.CreatePurchaseRequest true "Purchase request"
// @Success 201 {object} dtos.PurchaseResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/purchase [post]
func createPurchase(c *fiber.Ctx) error {
	return proxyToPurchaseService(c, "POST", "/api/v1/purchase")
}

// @Summary Upload payment proof for a purchase
// @Description Customer can upload their payment proof photo here. After payment, decreases the real product quantity.
// @Tags purchase
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param purchaseId path string true "Purchase ID"
// @Param request body dtos.PaymentProofRequest true "Payment proof request"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/purchase/{purchaseId} [post]
func uploadPaymentProof(c *fiber.Ctx) error {
	purchaseID := c.Params("purchaseId")
	return proxyToPurchaseService(c, "POST", "/api/v1/purchase/"+purchaseID)
}

// @Summary Get purchase by ID
// @Description Get a specific purchase by its ID
// @Tags purchase
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param purchaseId path string true "Purchase ID"
// @Success 200 {object} dtos.GetPurchaseResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/purchase/{purchaseId} [get]
func getPurchaseByID(c *fiber.Ctx) error {
	purchaseID := c.Params("purchaseId")
	return proxyToPurchaseService(c, "GET", "/api/v1/purchase/"+purchaseID)
}

// @Summary List user's purchases
// @Description Get a paginated list of user's purchases
// @Tags purchase
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} dtos.ListPurchasesResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/purchase [get]
func listPurchases(c *fiber.Ctx) error {
	return proxyToPurchaseService(c, "GET", "/api/v1/purchase")
}

// proxyToPurchaseService forwards requests to the purchase service
func proxyToPurchaseService(c *fiber.Ctx, method string, endpoint string) error {
	// Get request body
	body := c.Body()

	// Create request to purchase service
	purchase_service_url := c.Locals("service_urls").(*config.ServiceURLs).PurchaseServiceURL
	url := purchase_service_url + endpoint
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

	// Add user context headers (for protected routes)
	if userID := c.Locals("user_id"); userID != nil {
		req.Header.Set("X-User-ID", userID.(string))
		req.Header.Set("X-Auth-Gateway", "backend-infra")

		// Add internal secret for secure communication
		req.Header.Set("X-Secret", "backend-infra-internal-secret") // TODO: Make configurable
	}

	// Make request to purchase service
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
			"error": "Purchase service unavailable",
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
