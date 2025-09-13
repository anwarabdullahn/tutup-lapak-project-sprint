package config

import (
	"fmt"

	"github.com/spf13/viper"
)

// NewViper is a function to load config from .env file
func NewViper() *viper.Viper {
	config := viper.New()

	// Set config to read from .env file
	// Set config file name and type
	config.SetConfigName(".env")
	config.SetConfigType("env")
	config.AddConfigPath(".")
	// Read config file
	if err := config.ReadInConfig(); err != nil {
		fmt.Println("Warning: .env file not found or unreadable, using environment variables only:", err)
	} else {
		fmt.Println("Config file loaded from .env")
	}

	// Set config to read from environment variables (will override .env)
	config.AutomaticEnv()

	return config
}
