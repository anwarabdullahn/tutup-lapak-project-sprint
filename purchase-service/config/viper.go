package config

import (
	"github.com/spf13/viper"
)

// NewViper is a function to load config from .env file
func NewViper() *viper.Viper {
	config := viper.New()

	// Set config file name and type
	config.SetConfigName(".env")
	config.SetConfigType("env")
	config.AddConfigPath(".")
	
	// Read config file
	if err := config.ReadInConfig(); err != nil {
		// If .env file is not found, just use environment variables
	}
	
	// Set config to read from environment variables (will override .env)
	config.AutomaticEnv()

	return config
}