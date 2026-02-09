package server

import (
	"github.com/gofiber/fiber/v3"
)

// RegisterRoutes wires all HTTP routes for the application.
// Split into a separate file from server.go to keep routing concerns isolated.
func RegisterRoutes(app *fiber.App) {
	// Health and readiness endpoints
	app.Get("/healthz", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "ok",
		})
	})

	app.Get("/readyz", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// TODO: add domain-specific routes here, e.g.:
	// api := app.Group("/api")
	// v1 := api.Group("/v1")
	// v1.Get("/items", listItemsHandler)
}
