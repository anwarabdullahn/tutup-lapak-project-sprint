package main

// @title Product Service API
// @version 1.0
// @description API for managing products.
// @BasePath /api/v1

import (
	"log"
	"os"
	"product-service/api/handlers"
	"product-service/api/routes"
	"product-service/config"
	"product-service/pkg/entities"
	prod "product-service/pkg/product"
)

func main() {
	v := config.NewViper()
	db := config.NewGorm(v)

	// Auto-migrate ONLY for dev convenience
	if os.Getenv("AUTO_MIGRATE") == "true" {
		if err := db.AutoMigrate(&entities.Product{}); err != nil {
			log.Fatalf("auto migrate failed: %v", err)
		}
	}

	app := config.NewFiber()

	repo := prod.NewRepository(db)
	svc := prod.NewService(repo)
	h := handlers.NewProductHandler(svc)

	api := app.Group("/api/v1")
	routes.RegisterProductRoutes(api, h)

	port := v.GetString("PORT")
	if port == "" {
		port = "3003"
	}
	log.Printf("product-service listening on :%s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}
