package main

import (
	"github.com/anshu4sharma/resume_ats/internal/di"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			di.ConfigProvider,
			di.LoggerProvider,
			di.RedisProvider,
			di.AtsServiceProvider,
			di.DeployServiceProvider,
			di.HandlersProvider,
			di.FiberApp,
		),
		fx.Invoke(
			di.RegisterHooks,
		),
	).Run()
}
