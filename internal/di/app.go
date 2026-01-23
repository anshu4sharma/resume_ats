package di

import (
	"context"

	"github.com/anshu4sharma/resume_ats/internal/config"
	"github.com/anshu4sharma/resume_ats/internal/handlers"
	"github.com/anshu4sharma/resume_ats/internal/router"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
)

func FiberApp() *fiber.App {
	app := fiber.New()

	app.Static("/", "./web") 

	return app
}

func RegisterHooks(
	lc fx.Lifecycle,
	app *fiber.App,
	cfg *config.Config,
	handlers *handlers.Handlers,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			router.SetupRoutes(app, handlers)
			go app.Listen(cfg.ServerPort)
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return app.Shutdown()
		},
	})
}
