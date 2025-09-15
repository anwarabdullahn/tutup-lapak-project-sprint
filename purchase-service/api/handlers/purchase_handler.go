package handlers

import (
	"purchase-service/pkg/dtos"
	"purchase-service/pkg/purchase"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type PurchaseHandler struct {
	service   purchase.Service
	validator *validator.Validate
}

func NewPurchaseHandler(service purchase.Service) *PurchaseHandler {
	return &PurchaseHandler{
		service:   service,
		validator: validator.New(),
	}
}

// CreatePurchase handles POST /v1/purchase
// @Summary Create a new purchase
// @Description Customer can add their items to cart so they can pay them
// @Tags purchase
// @Accept json
// @Produce json
// @Param request body dtos.CreatePurchaseRequest true "Purchase request"
// @Success 201 {object} presenter.PurchaseResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/purchase [post]
func (h *PurchaseHandler) CreatePurchase(c *fiber.Ctx) error {
	var req dtos.CreatePurchaseRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		validationErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors[err.Field()] = getValidationMessage(err)
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": validationErrors,
		})
	}

	// Additional validation for contact details
	if err := h.validateContactDetails(req.SenderContactType, req.SenderContactDetail); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Create purchase
	purchase, err := h.service.CreatePurchase(c.Context(), req)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create purchase",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(purchase)
}

// UploadPaymentProof handles POST /v1/purchase/:purchaseId
// @Summary Upload payment proof for a purchase
// @Description Customer can upload their payment proof photo here. After payment, decreases the real product quantity.
// @Tags purchase
// @Accept json
// @Produce json
// @Param purchaseId path string true "Purchase ID"
// @Param request body dtos.PaymentProofRequest true "Payment proof request"
// @Success 201 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/purchase/{purchaseId} [post]
func (h *PurchaseHandler) UploadPaymentProof(c *fiber.Ctx) error {
	purchaseID := c.Params("purchaseId")
	if purchaseID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Purchase ID is required",
		})
	}

	var req dtos.PaymentProofRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		validationErrors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			validationErrors[err.Field()] = getValidationMessage(err)
		}
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": validationErrors,
		})
	}

	// Upload payment proof
	if err := h.service.UploadPaymentProof(c.Context(), purchaseID, req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload payment proof",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Payment proof uploaded successfully",
	})
}

// GetPurchaseByID handles GET /v1/purchase/:purchaseId
// @Summary Get purchase by ID
// @Description Get a specific purchase by its ID
// @Tags purchase
// @Accept json
// @Produce json
// @Param purchaseId path string true "Purchase ID"
// @Success 200 {object} presenter.GetPurchaseResponse
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/purchase/{purchaseId} [get]
func (h *PurchaseHandler) GetPurchaseByID(c *fiber.Ctx) error {
	purchaseID := c.Params("purchaseId")
	if purchaseID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Purchase ID is required",
		})
	}

	// Get purchase
	purchase, err := h.service.GetPurchaseByID(c.Context(), purchaseID)
	if err != nil {
		if err.Error() == "purchase not found" || err.Error() == "invalid purchase ID format" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Purchase not found",
			})
		}
		if err.Error() == "unauthorized: purchase does not belong to user" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get purchase",
		})
	}

	return c.Status(fiber.StatusOK).JSON(purchase)
}

// ListPurchases handles GET /v1/purchase
// @Summary List user's purchases
// @Description Get a paginated list of user's purchases
// @Tags purchase
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} presenter.ListPurchasesResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /v1/purchase [get]
func (h *PurchaseHandler) ListPurchases(c *fiber.Ctx) error {
	// Parse query parameters
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	// Validate pagination parameters
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// Get purchases
	purchases, err := h.service.ListPurchases(c.Context(), page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get purchases",
		})
	}

	return c.Status(fiber.StatusOK).JSON(purchases)
}

// validateContactDetails validates email or phone based on contact type
func (h *PurchaseHandler) validateContactDetails(contactType, contactDetail string) error {
	if contactType == "email" {
		if err := h.validator.Var(contactDetail, "email"); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid email format")
		}
	} else if contactType == "phone" {
		if err := h.validator.Var(contactDetail, "min=10,max=15"); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid phone number format")
		}
	}
	return nil
}

// getValidationMessage returns a user-friendly validation message
func getValidationMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "This field is required"
	case "min":
		return "Value is too short"
	case "max":
		return "Value is too long"
	case "email":
		return "Invalid email format"
	case "oneof":
		return "Invalid value, must be one of: " + err.Param()
	case "dive":
		return "Invalid array item"
	default:
		return "Invalid value"
	}
}
