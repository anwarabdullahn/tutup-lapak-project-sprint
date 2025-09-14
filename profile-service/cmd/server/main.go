package main

import (
	"log"
	"profile-service/api/routes"
	"profile-service/config"
	"profile-service/pkg/auth"
	"profile-service/pkg/entities"
	"profile-service/pkg/user"
	"profile-service/pkg/userfile"
	"time"

	"github.com/gofiber/contrib/swagger"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// @title My API
// @version 1.0
// @description This is my API with JWT auth

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Load config.yaml pake Viper
	v := config.NewViper()

	// // Init DB (GORM)
	db := config.NewGorm(v)

	// // Auto-migrate the entities
	err := db.AutoMigrate(&entities.File{}, &entities.User{})

	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}

	// Init Auth Service & Repo
	// authRepo := auth.NewGormRepository(db)
	// authService := auth.NewService(authRepo)

	profileRepo := user.NewGormRepository(db)
	profileService := user.NewService(profileRepo)

	uploadFileRepo := userfile.NewRepo(db)
	uploadFileService := userfile.NewService(uploadFileRepo)

	app := config.NewFiber(v)

	app.Use(cors.New())

	// Init JWT Manager (24 jam expired)
	jwtManager := auth.NewJWTManager(v.GetString("jwt.secret"), 24*time.Hour)

	// Init Fiber

	api := app.Group("/api/v1")
	// routes.AuthRouter(api, authService, jwtManager)
	routes.ProfileRouter(api, profileService, jwtManager)
	routes.UserfileRouter(api, uploadFileService, jwtManager)

	app.Use(swagger.New(swagger.Config{
		BasePath: "/",
		FilePath: "./docs/swagger.json",
		Path:     "swagger",
		Title:    "Swagger API Docs",
		CacheAge: 86400,
	}))

	// Run server
	port := v.GetString("server.port")
	if port == "" {
		port = "3002"
	}
	log.Fatal(app.Listen(":" + port))
}
