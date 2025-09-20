package routes

import (
	"product-service/api/handlers"

	"github.com/gofiber/fiber/v2"
)

func RegisterProductRoutes(r fiber.Router, h *handlers.ProductHandler) {
	g := r.Group("/products")
	g.Post("/", h.Create)
	g.Get("/", h.List)
	g.Get(":id", h.Get)
	g.Patch(":id", h.Update)
	g.Delete(":id", h.Delete)
}
