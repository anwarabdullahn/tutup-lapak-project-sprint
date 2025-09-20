package config

import (
	"github.com/gofiber/fiber/v2"
)

func NewFiber() *fiber.App {
	return fiber.New(fiber.Config{AppName: "product-service"})
}
