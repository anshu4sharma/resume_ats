package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/anshu4sharma/resume_ats/internal/bootstrap"
	"github.com/anshu4sharma/resume_ats/internal/config"
	"github.com/anshu4sharma/resume_ats/internal/router"
	"github.com/anshu4sharma/resume_ats/pkg/utils"
	"github.com/gofiber/fiber/v2"
	flogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	cfg := config.Load()

	logger := utils.NewLogger()

	app := fiber.New(fiber.Config{
		AppName:   "Resume ATS v1.0",
		BodyLimit: utils.MaxResumeSizeBytes,
	})

	app.Use(recover.New())
	app.Use(flogger.New())
	app.Static("/", "./web")
	
	handlers, err := bootstrap.InitializeApp(logger)
	
	if err != nil {
		logger.Errorf("Failed to initialize app: %v", err)
		os.Exit(1)
	}

	router.SetupRoutes(app, handlers)

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-shutdown
		logger.Infof("Shutting down server...")
		if err := app.Shutdown(); err != nil {
			logger.Errorf("Error shutting down server: %v", err)
		}
	}()

	log.Fatal(app.Listen(cfg.ServerPort))
}
