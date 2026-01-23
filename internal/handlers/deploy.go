package handlers

import (
	"github.com/anshu4sharma/resume_ats/internal/config"
	"github.com/anshu4sharma/resume_ats/internal/services"
	"github.com/anshu4sharma/resume_ats/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type DeployHandler struct {
	deployService services.DeployService
	cfg           *config.Config
}

func NewDeployHandler(
	deployService services.DeployService,
	logger *utils.Logger,
	cfg *config.Config,
) *DeployHandler {
	return &DeployHandler{
		deployService: deployService,
		cfg:           cfg,
	}
}

func (h *DeployHandler) Deploy(c *fiber.Ctx) error {
	sig := c.Get("X-Hub-Signature-256")
	if sig == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	if !utils.VerifyGitHubSignature(
		h.cfg.DEPLOY_SECRET,
		c.Body(),
		sig,
	) {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	h.deployService.TriggerDeploy(c.Context())

	return c.JSON(fiber.Map{
		"status": "deploy triggered",
	})
}
