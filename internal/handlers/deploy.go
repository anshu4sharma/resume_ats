package handlers

import (
	"os"
	"os/exec"

	"github.com/anshu4sharma/resume_ats/internal/config"
	"github.com/anshu4sharma/resume_ats/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type DeployHandler struct {
	logger *utils.Logger
}

func NewDeployHandler(logger *utils.Logger) *DeployHandler {
	return &DeployHandler{logger: logger}
}

func (h *DeployHandler) Deploy(c *fiber.Ctx) error {
	sig := c.Get("X-Hub-Signature-256")
	if sig == "" {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	body := c.Body()

	if !utils.VerifyGitHubSignature(
		config.Load().DEPLOY_SECRET,
		body,
		sig,
	) {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	go func() {
		cmd := exec.Command("/bin/bash", "/home/ubuntu/resume_ats/deploy.sh")
		cmd.Dir = "/home/ubuntu/resume_ats"
		cmd.Env = os.Environ()

		if out, err := cmd.CombinedOutput(); err != nil {
			h.logger.Errorf("deploy failed: %v | output: %s", err, out)
		} else {
			h.logger.Infof("deploy succeeded: %s", out)
		}
	}()

	return c.JSON(fiber.Map{
		"status": "deploy triggered",
	})
}
