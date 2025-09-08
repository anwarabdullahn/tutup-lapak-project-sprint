package config

import (
	"auth-service/pkg/user"

	"gorm.io/gorm"
)

// Services struct holds all service dependencies
type Services struct {
	UserService user.Service
}

// InitServices initializes all application services
func InitServices(db *gorm.DB) Services {
	// Initialize repositories
	userRepo := user.NewGormRepository(db)

	// Initialize services
	userService := user.NewService(userRepo)

	return Services{
		UserService: userService,
	}
}
