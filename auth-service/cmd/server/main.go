package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	svc := getEnv("SERVICE_NAME", "auth-service")
	port := getEnv("PORT", "8080")

	app := fiber.New()

	app.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"service": svc,
			"status":  "running",
		})
	})

	addr := ":" + port
	log.Printf("%s listening on %s", svc, addr)
	if err := app.Listen(addr); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
