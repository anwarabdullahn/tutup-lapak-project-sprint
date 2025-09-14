package config

import (
	"github.com/spf13/viper"
)

type ServiceURLs struct {
	AuthServiceURL      string
	PROFILE_SERVICE_URL string
	// Add more as needed
}

func LoadServiceURLs(config *viper.Viper) *ServiceURLs {
	AuthServiceURL := config.GetString("AUTH_SERVICE_URL")
	PROFILE_SERVICE_URL := config.GetString("PROFILE_SERVICE_URL")
	if AuthServiceURL == "" {
		AuthServiceURL = "http://localhost:3001"
	}

	if PROFILE_SERVICE_URL == "" {
		PROFILE_SERVICE_URL = "http://localhost:3002"
	}

	return &ServiceURLs{
		AuthServiceURL:      AuthServiceURL,
		PROFILE_SERVICE_URL: PROFILE_SERVICE_URL,
	}
}
