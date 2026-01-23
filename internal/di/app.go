package di

import (
	"context"

	"github.com/anshu4sharma/resume_ats/internal/config"
	"github.com/anshu4sharma/resume_ats/internal/handlers"
	"github.com/anshu4sharma/resume_ats/internal/router"
	"github.com/anshu4sharma/resume_ats/pkg/utils"
	"github.com/gofiber/fiber/v2"
	flogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"go.uber.org/fx"
)

func FiberApp() *fiber.App {
	app := fiber.New(fiber.Config{
		AppName:   "Resume ATS v1.0",
		BodyLimit: utils.MaxResumeSizeBytes,
		Prefork:   false,
	})

	app.Use(recover.New())
	app.Use(flogger.New())
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

			go func() {
				if err := app.Listen(cfg.SERVER_PORT); err != nil {
					// Fiber returns an error on shutdown; ignore it
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			done := make(chan error, 1)

			go func() {
				done <- app.Shutdown()
			}()

			select {
			case <-ctx.Done():
				return ctx.Err()
			case err := <-done:
				return err
			}
		},
	})
}
