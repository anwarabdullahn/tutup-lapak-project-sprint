package main

import (
    "fmt"
    "log"
    "os"

    "github.com/gofiber/fiber/v2"
    "github.com/gofiber/fiber/v2/middleware/logger"
    "github.com/gofiber/fiber/v2/middleware/proxy"
)

func main() {
    port := getEnv("PORT", "3000")

    app := fiber.New()
    app.Use(logger.New())

    // Healthcheck
    app.Get("/healthz", func(c *fiber.Ctx) error { return c.SendString("ok") })

    // Proxy helpers
    proxyTo := func(base string) fiber.Handler {
        return func(c *fiber.Ctx) error {
            // Strip the first path segment (e.g., /auth or /profile)
            rest := c.Params("*")
            target := base
            if rest != "" {
                target = fmt.Sprintf("%s/%s", base, rest)
            }
            return proxy.Do(c, target)
        }
    }

    // Routes
    app.All("/auth", proxyTo("http://auth-service:3001"))
    app.All("/auth/*", proxyTo("http://auth-service:3001"))

    app.All("/profile", proxyTo("http://profile-service:3002"))
    app.All("/profile/*", proxyTo("http://profile-service:3002"))

    addr := ":" + port
    log.Printf("gateway listening on %s", addr)
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

