package bootstrap

import (
	"github.com/anshu4sharma/resume_ats/internal/handlers"
	"github.com/anshu4sharma/resume_ats/internal/services"
	"github.com/anshu4sharma/resume_ats/pkg/utils"
)

func InitializeApp(logger *utils.Logger) (*handlers.Handlers, error) {
	atsService := services.NewAtsService(logger)
	handlers := handlers.NewHandlers(
		atsService,
		logger,
	)
	return handlers, nil
}
