package routes

import (
	"backend-infra/config"
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func SetupProductRoutes(app *fiber.App) {
	g := app.Group("/v1/products")
	g.Post("/", createProduct)
	g.Get("/", listProducts)
	g.Get(":id", getProduct)
	g.Patch(":id", updateProduct)
	g.Delete(":id", deleteProduct)
}

// @Summary Create product
// @Tags products
// @Accept json
// @Produce json
// @Param request body dtos.CreateProductRequest true "Product data"
// @Success 201 {object} dtos.Product
// @Failure 400 {object} map[string]any
// @Router /v1/products/ [post]
func createProduct(c *fiber.Ctx) error { return proxyToProductService(c, "POST", "/api/v1/products/") }

// @Summary List products
// @Tags products
// @Produce json
// @Success 200 {array} dtos.Product
// @Router /v1/products/ [get]
func listProducts(c *fiber.Ctx) error { return proxyToProductService(c, "GET", "/api/v1/products/") }

// @Summary Get product
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dtos.Product
// @Failure 404 {object} map[string]any
// @Router /v1/products/{id} [get]
func getProduct(c *fiber.Ctx) error {
	return proxyToProductService(c, "GET", "/api/v1/products/"+c.Params("id"))
}

// @Summary Update product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param request body dtos.UpdateProductRequest true "Product data"
// @Success 200 {object} dtos.Product
// @Failure 400 {object} map[string]any
// @Router /v1/products/{id} [patch]
func updateProduct(c *fiber.Ctx) error {
	return proxyToProductService(c, "PATCH", "/api/v1/products/"+c.Params("id"))
}

// @Summary Delete product
// @Tags products
// @Param id path string true "Product ID"
// @Success 204
// @Failure 404 {object} map[string]any
// @Router /v1/products/{id} [delete]
func deleteProduct(c *fiber.Ctx) error {
	return proxyToProductService(c, "DELETE", "/api/v1/products/"+c.Params("id"))
}

func proxyToProductService(c *fiber.Ctx, method, endpoint string) error {
	body := c.Body()
	urls := c.Locals("service_urls").(*config.ServiceURLs)
	url := urls.ProductServiceURL + endpoint
	req, err := http.NewRequest(method, url, bytes.NewReader(body))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create request"})
	}
	req.Header.Set("Content-Type", "application/json")
	for k, vs := range c.GetReqHeaders() {
		for _, v := range vs {
			req.Header.Add(k, v)
		}
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{"error": "Product service unavailable"})
	}
	defer resp.Body.Close()
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read response"})
	}
	c.Status(resp.StatusCode)
	for k, vs := range resp.Header {
		for _, v := range vs {
			c.Set(k, v)
		}
	}
	var jsonResp any
	if err := json.Unmarshal(respBody, &jsonResp); err != nil {
		return c.SendString(string(respBody))
	}
	return c.JSON(jsonResp)
}
