package config

import (
	"profile-service/pkg/uploadfile"
	"profile-service/pkg/user"

	"gorm.io/gorm"
)

// Services struct holds all service dependencies
type Services struct {
	UserService user.Service
	FileService uploadfile.Service
}

// InitServices initializes all application services
func InitServices(db *gorm.DB) Services {
	// Initialize repositories
	userRepo := user.NewGormRepository(db)

	// Initialize services
	userService := user.NewService(userRepo)

	fileRepo := uploadfile.NewRepo(db)
	fileService := uploadfile.NewService(fileRepo)

	// userfile.Service

	return Services{
		UserService: userService,
		FileService: fileService,
	}
}
