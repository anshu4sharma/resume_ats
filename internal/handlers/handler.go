package handlers

import (
	"github.com/anshu4sharma/resume_ats/internal/config"
	"github.com/anshu4sharma/resume_ats/internal/services"
	redis "github.com/anshu4sharma/resume_ats/pkg"
	"github.com/anshu4sharma/resume_ats/pkg/utils"
)

type Handlers struct {
	AtsHandler   *AtsHandler
	DeployHanler *DeployHandler
}

func NewHandlers(atsService services.AtsService, deploy services.DeployService, logger *utils.Logger, cfg *config.Config, redis *redis.RedisClient) *Handlers {
	return &Handlers{
		AtsHandler:   NewAtsHandler(atsService, logger, cfg, redis),
		DeployHanler: NewDeployHandler(deploy, logger, cfg),
	}
}
