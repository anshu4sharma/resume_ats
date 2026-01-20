package handlers

import (
	"os"

	"github.com/anshu4sharma/resume_ats/internal/services"
	"github.com/anshu4sharma/resume_ats/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type AtsHandler struct {
	service *services.AtsService
	logger  *utils.Logger
}

func NewAtsHandler(service *services.AtsService, logger *utils.Logger) *AtsHandler {
	return &AtsHandler{service: service, logger: logger}
}

func (h *AtsHandler) UploadResume(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			SendString("file is required")
	}

	if file.Header.Get("Content-Type") != "application/pdf" {
		return c.Status(fiber.StatusBadRequest).
			SendString("only PDF files are supported")
	}

	uploadDir := "./uploads"
	_ = os.MkdirAll(uploadDir, 0755)

	path := uploadDir + "/" + file.Filename
	if err := c.SaveFile(file, path); err != nil {
		return c.Status(500).SendString("failed to save file")
	}
	
	if file.Size > utils.MaxResumeSizeBytes {
		return fiber.ErrRequestEntityTooLarge
	}

	defer func() {
		h.logger.Debugf("Cleaning up resume pdf.")
		if err := os.Remove(path); err != nil {
			h.logger.Warnf("failed to cleanup temp file %s: %v", path, err)
		}
	}()

	result, err := h.service.AnalyzeResume(
		path,
		file.Size,
		file.Filename,
	)
	if err != nil {
		return c.Status(400).SendString(err.Error())
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"score":  result.Score,
		"data":   result.Data,
	})
}
