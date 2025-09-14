package config

import (
	"github.com/spf13/viper"
)

type ServiceURLs struct {
	AuthServiceURL     string
	PurchaseServiceURL string
	// Add more as needed
}

func LoadServiceURLs(config *viper.Viper) *ServiceURLs {
	AuthServiceURL := config.GetString("AUTH_SERVICE_URL")
	if AuthServiceURL == "" {
		AuthServiceURL = "http://localhost:3001"
	}

	PurchaseServiceURL := config.GetString("PURCHASE_SERVICE_URL")
	if PurchaseServiceURL == "" {
		PurchaseServiceURL = "http://localhost:3004"
	}

	return &ServiceURLs{
		AuthServiceURL:     AuthServiceURL,
		PurchaseServiceURL: PurchaseServiceURL,
	}
}
