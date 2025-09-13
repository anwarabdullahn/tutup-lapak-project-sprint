package config

import (
	"github.com/spf13/viper"
)

type ServiceURLs struct {
	AuthServiceURL string
	// Add more as needed
}

func LoadServiceURLs(config *viper.Viper) *ServiceURLs {
	AuthServiceURL := config.GetString("AUTH_SERVICE_URL")
	if AuthServiceURL == "" {
		AuthServiceURL = "http://localhost:3001"
	}

	return &ServiceURLs{
		AuthServiceURL: AuthServiceURL,
	}
}
