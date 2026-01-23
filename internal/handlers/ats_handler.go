package handlers

import (
	"encoding/json"
	"os"
	"time"

	"github.com/anshu4sharma/resume_ats/internal/config"
	"github.com/anshu4sharma/resume_ats/internal/services"
	"github.com/anshu4sharma/resume_ats/internal/structs"
	redis "github.com/anshu4sharma/resume_ats/pkg"
	"github.com/anshu4sharma/resume_ats/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type AtsHandler struct {
	service services.AtsService
	logger  *utils.Logger
	cfg     *config.Config
	redis   *redis.RedisClient
}

func NewAtsHandler(
	service services.AtsService,
	logger *utils.Logger,
	cfg *config.Config,
	redis *redis.RedisClient,
) *AtsHandler {
	return &AtsHandler{
		service: service,
		logger:  logger,
		cfg:     cfg,
		redis:   redis,
	}
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

	fileHash, err := utils.HashFile(path)

	if err != nil {
		return c.JSON(fiber.Map{
			"status": "failed",
			"error":  err.Error(),
		})
	}

	cached, err := h.redis.GetValue(c.Context(), fileHash)
	if err == nil && cached != "" {
		var cachedResult structs.ResumeAnalysisResult
		if err := json.Unmarshal([]byte(cached), &cachedResult); err == nil {
			h.logger.Infof("cache hit for %s", fileHash)

			return c.JSON(fiber.Map{
				"status": "success",
				"score":  cachedResult.Score,
				"data":   cachedResult.Data,
			})
		}

		h.logger.Warnf("cache unmarshal failed, recomputing")
	}

	result, err := h.service.AnalyzeResume(
		path,
		file.Size,
		file.Filename,
	)

	if err != nil {
		return c.JSON(fiber.Map{
			"status": "failed",
			"error":  err.Error(),
		})
	}

	bytes, err := json.Marshal(result)
	if err == nil {
		_ = h.redis.SetValue(
			c.Context(),
			fileHash,
			string(bytes),
			time.Hour,
		)
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"score":  result.Score,
		"data":   result.Data,
	})
}
