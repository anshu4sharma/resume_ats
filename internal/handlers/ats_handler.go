package handlers

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	"github.com/anshu4sharma/resume_ats/internal/config"
	"github.com/anshu4sharma/resume_ats/internal/services"
	"github.com/anshu4sharma/resume_ats/internal/structs"
	redis "github.com/anshu4sharma/resume_ats/pkg"
	"github.com/anshu4sharma/resume_ats/pkg/utils"
	"github.com/anshu4sharma/resume_ats/shared/constant"
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

	if file.Size > utils.MaxResumeSize {
		return c.Status(fiber.StatusBadRequest).
			SendString("file size exceeds 3 MB.")
	}

	uploadDir := "./uploads"

	_ = os.MkdirAll(uploadDir, 0755)

	path := uploadDir + "/" + utils.GenerateID(4) + "-" + file.Filename

	src, err := file.Open()

	if err != nil {
		h.logger.Infof("error while opening file in memory")
	}

	defer src.Close()

	data, err := io.ReadAll(src)
	if err != nil {
		return c.Status(500).SendString("failed to read file")
	}

	fileHash := utils.HashBytes(data)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "failed",
			"error":  "failed to hash file",
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

	if err := c.SaveFile(file, path); err != nil {
		return c.Status(500).SendString("failed to save file")
	}

	reader := bytes.NewReader(data)

	result, err := h.service.AnalyzeResume(
		reader,
		file.Filename,
	)

	if err != nil {
		return c.Status(500).JSON(fiber.Map{
			"status": "failed",
			"error":  "Sorry we couldnt anaylyze ur resume",
		})
	}

	bytes, err := json.Marshal(result)
	if err == nil {
		_ = h.redis.SetValue(
			c.Context(),
			fileHash,
			string(bytes),
			constant.Max_Cache_Duration,
		)
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"score":  result.Score,
		"data":   result.Data,
	})
}
