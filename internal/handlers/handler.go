package handlers

import (
	"github.com/anshu4sharma/resume_ats/internal/services"
	"github.com/anshu4sharma/resume_ats/pkg/utils"
)

type Handlers struct {
	AtsHandler *AtsHandler
}

func NewHandlers(atsService *services.AtsService, logger *utils.Logger) *Handlers {
	return &Handlers{
		AtsHandler: NewAtsHandler(atsService, logger),
	}
}
