package handlers

import (
	"os/exec"

	"github.com/anshu4sharma/resume_ats/internal/config"
	"github.com/gofiber/fiber/v2"
)

func (h *Handlers) Deploy(c *fiber.Ctx) error {
	secret := c.Get("X-Deploy-Secret")
	if secret != config.Load().DEPLOY_SECRET {
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	go exec.Command("/home/ubuntu/resume_ats/deploy.sh").Run()

	return c.JSON(fiber.Map{
		"status": "deploy triggered",
	})
}
