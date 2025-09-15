package config

import (
	"os"
	"purchase-service/pkg/purchase"

	"gorm.io/gorm"
)

// Services struct holds all service dependencies
type Services struct {
	PurchaseService purchase.Service
}

// InitServices initializes all application services
func InitServices(db *gorm.DB) Services {
	// Initialize repositories
	purchaseRepo := purchase.NewGormRepository(db)

	// Get service URLs from environment variables
	userServiceURL := os.Getenv("USER_SERVICE_URL")
	if userServiceURL == "" {
		userServiceURL = "http://localhost:3002" // Default user service URL
	}

	productServiceURL := os.Getenv("PRODUCT_SERVICE_URL")
	if productServiceURL == "" {
		productServiceURL = "http://localhost:3003" // Default product service URL
	}

	internalSecret := os.Getenv("INTERNAL_SECRET")
	if internalSecret == "" {
		internalSecret = "backend-infra-internal-secret" // Default internal secret
	}

	// Initialize services
	purchaseService := purchase.NewService(purchaseRepo, userServiceURL, productServiceURL, internalSecret)

	return Services{
		PurchaseService: purchaseService,
	}
}
