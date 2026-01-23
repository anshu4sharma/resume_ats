package router

import (
	"github.com/anshu4sharma/resume_ats/internal/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, h *handlers.Handlers) {
	apiV1 := app.Group("/api/v1")

	ats := apiV1.Group("/ats")
	ats.Post("/resume_upload", h.AtsHandler.UploadResume)

	deploy := apiV1.Group("/deploy")
	deploy.Post("/webhook", h.Deploy)
}
