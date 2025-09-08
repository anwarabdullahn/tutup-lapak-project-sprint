package config

import (
	"auth-service/api/routes"
	"auth-service/pkg/user"

	"gorm.io/gorm"
)

// InitServices initializes all application services
func InitServices(db *gorm.DB) routes.Services {
	// Initialize repositories
	userRepo := user.NewGormRepository(db)

	// Initialize services
	userService := user.NewService(userRepo)

	return routes.Services{
		UserService: userService,
	}
}
