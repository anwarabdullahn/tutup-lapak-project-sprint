package handlers

import (
	"net/http"
	"product-service/pkg/product"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ProductHandler struct {
	svc product.Service
}

func NewProductHandler(svc product.Service) *ProductHandler { return &ProductHandler{svc: svc} }

// CreateProductRequest represents create payload
// @description Create product request body
type CreateProductRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
	PriceCents  int64  `json:"price_cents"`
}

// UpdateProductRequest represents partial update payload
type UpdateProductRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
	PriceCents  *int64  `json:"price_cents"`
}

// Create product
// @Summary Create product
// @Tags products
// @Accept json
// @Produce json
// @Param data body CreateProductRequest true "product data"
// @Success 201 {object} map[string]any
// @Failure 400 {object} fiber.Error
// @Router /products/ [post]
func (h *ProductHandler) Create(c *fiber.Ctx) error {
	var req CreateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid json body")
	}
	input := product.CreateInput{Name: req.Name, Description: req.Description, Category: req.Category, PriceCents: req.PriceCents}
	p, err := h.svc.Create(c.Context(), input)
	if err != nil {
		return errorResponse(err)
	}
	return c.Status(http.StatusCreated).JSON(p)
}

// Get product by ID
// @Summary Get product
// @Tags products
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} map[string]any
// @Failure 404 {object} fiber.Error
// @Router /products/{id} [get]
func (h *ProductHandler) Get(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid id")
	}
	p, err := h.svc.Get(c.Context(), id)
	if err != nil {
		return errorResponse(err)
	}
	if p == nil {
		return fiber.NewError(http.StatusNotFound, "not found")
	}
	return c.JSON(p)
}

// List products
// @Summary List products
// @Tags products
// @Produce json
// @Success 200 {array} map[string]any
// @Router /products/ [get]
func (h *ProductHandler) List(c *fiber.Ctx) error {
	ps, err := h.svc.List(c.Context())
	if err != nil {
		return errorResponse(err)
	}
	return c.JSON(ps)
}

// Update product partially
// @Summary Update product
// @Tags products
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Param data body UpdateProductRequest true "update data"
// @Success 200 {object} map[string]any
// @Failure 400 {object} fiber.Error
// @Router /products/{id} [patch]
func (h *ProductHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid id")
	}
	var req UpdateProductRequest
	if err := c.BodyParser(&req); err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid json body")
	}
	input := product.UpdateInput{Name: req.Name, Description: req.Description, Category: req.Category, PriceCents: req.PriceCents}
	p, err := h.svc.Update(c.Context(), id, input)
	if err != nil {
		return errorResponse(err)
	}
	return c.JSON(p)
}

// Delete product
// @Summary Delete product
// @Tags products
// @Param id path string true "Product ID"
// @Success 204
// @Failure 404 {object} fiber.Error
// @Router /products/{id} [delete]
func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return fiber.NewError(http.StatusBadRequest, "invalid id")
	}
	if err := h.svc.Delete(c.Context(), id); err != nil {
		return errorResponse(err)
	}
	return c.SendStatus(http.StatusNoContent)
}

func errorResponse(err error) *fiber.Error {
	switch err.Error() {
	case "invalid category":
		return fiber.NewError(http.StatusBadRequest, err.Error())
	case "product not found":
		return fiber.NewError(http.StatusNotFound, err.Error())
	default:
		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}
}
