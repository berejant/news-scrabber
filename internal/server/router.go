package server

import (
	"news-scrabber/internal/server/actions/transribe"

	"github.com/gofiber/fiber/v3"
)

// RegisterRoutes wires all HTTP routes for the application.
// Split into a separate file from server.go to keep routing concerns isolated.
func RegisterRoutes(app *fiber.App, act *transribe.RequestTranscribeAction) {
	// Health and readiness endpoints
	app.Get("/healthz", func(c fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status": "ok",
		})
	})

	app.Get("/readyz", func(c fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	// Transcription API
	v1 := app.Group("/api/v1")
	v1.Post("/transcribe-requests", act.Handle)
}
