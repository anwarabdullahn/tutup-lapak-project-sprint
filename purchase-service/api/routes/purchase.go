package routes

import (
	"purchase-service/api/handlers"
	"purchase-service/api/middleware"
	"purchase-service/config"

	"github.com/gofiber/fiber/v2"
)

// PurchaseRouter sets up purchase-related routes
func PurchaseRouter(api fiber.Router, services config.Services) {
	purchaseHandler := handlers.NewPurchaseHandler(services.PurchaseService)

	config := config.NewViper()

	// Purchase routes
	purchase := api.Group("/purchase")
	{
		purchase.Post("/", middleware.GatewayTrust(config), purchaseHandler.CreatePurchase)
		purchase.Get("/", middleware.GatewayTrust(config), purchaseHandler.ListPurchases)
		purchase.Get("/:purchaseId", middleware.GatewayTrust(config), purchaseHandler.GetPurchaseByID)
		purchase.Post("/:purchaseId", middleware.GatewayTrust(config), purchaseHandler.UploadPaymentProof)
	}
}
