package transribe

import (
	"news-scrabber/internal/transcribe"

	"github.com/gofiber/fiber/v3"
)

// RequestTranscribeAction implements the HTTP action for creating a transcribe request.
// Pattern: Actions — keep endpoint logic in small, testable units.
//
// POST /api/v1/transcribe-requests
// Body: {"url": "...", "job_id": "optional"}
// Returns: 202 {"job_id": "..."}
//
// It leverages TranscribeEventPublisher to emit an event to NATS.

// RequestTranscribeRequest is the expected JSON payload for the endpoint.
type RequestTranscribeRequest struct {
	URL   string `json:"url"`
	JobID string `json:"job_id"`
}

type RequestTranscribeAction struct {
	pub transcribe.TranscribeEventPublisher
}

func NewRequestTranscribeAction(pub transcribe.TranscribeEventPublisher) *RequestTranscribeAction {
	return &RequestTranscribeAction{pub: pub}
}

// Handle processes the HTTP request and publishes the corresponding event.
func (a *RequestTranscribeAction) Handle(c fiber.Ctx) error {
	var req RequestTranscribeRequest
	if err := c.Bind().JSON(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json"})
	}
	if req.URL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "url is required"})
	}
	jobID, err := a.pub.PublishVideoTranscribeRequested(c.Context(), req.URL, req.JobID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{"job_id": jobID})
}
