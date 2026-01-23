package services

import (
	"context"
	"os"
	"os/exec"
	"time"

	"github.com/anshu4sharma/resume_ats/internal/config"
	"github.com/anshu4sharma/resume_ats/pkg/utils"
)

type DeployService interface {
	TriggerDeploy(ctx context.Context)
}

type deployService struct {
	logger *utils.Logger
	config *config.Config
}

func NewDeployService(logger *utils.Logger, cf *config.Config) DeployService {
	return &deployService{logger: logger, config: cf}
}

func (s *deployService) TriggerDeploy(ctx context.Context) {
	go func() {
		ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
		defer cancel()

		cmd := exec.CommandContext(
			ctx,
			"/bin/bash",
			"/home/ubuntu/resume_ats/deploy.sh",
		)

		cmd.Dir = "/home/ubuntu/resume_ats"
		cmd.Env = os.Environ()

		out, err := cmd.CombinedOutput()
		if err != nil {
			s.logger.Errorf("deploy failed: %v | output: %s", err, out)
			return
		}

		s.logger.Infof("deploy succeeded: %s", out)
	}()
}
