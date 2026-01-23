package di

import (
	"context"

	"github.com/anshu4sharma/resume_ats/internal/config"
	"github.com/anshu4sharma/resume_ats/internal/handlers"
	"github.com/anshu4sharma/resume_ats/internal/services"
	redis "github.com/anshu4sharma/resume_ats/pkg"
	"github.com/anshu4sharma/resume_ats/pkg/utils"
	"github.com/anshu4sharma/resume_ats/shared/constant"
	"go.uber.org/fx"
)

func ConfigProvider() *config.Config {
	return config.Load()
}

func LoggerProvider() *utils.Logger {
	return utils.NewLogger()
}

func RedisProvider(
	lc fx.Lifecycle,
	cfg *config.Config,
	logger *utils.Logger,
) (*redis.RedisClient, error) {

	client := redis.NewRedisClient(cfg.REDIS_URL, logger)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return client.Connect(constant.RetryLimit, constant.RetryLimit)
		},
		OnStop: func(ctx context.Context) error {
			client.Close()
			return nil
		},
	})

	return client, nil
}

func AtsServiceProvider(lc fx.Lifecycle, logger *utils.Logger) services.AtsService {
	svc := services.NewAtsService(logger)
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Infof("ats service starting...")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Infof("ats service stopping")
			return nil
		},
	})

	return svc
}

func DeployServiceProvider(
	lc fx.Lifecycle,
	logger *utils.Logger,
	cf *config.Config,
) services.DeployService {

	svc := services.NewDeployService(logger, cf)

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			logger.Infof("deploy service starting...")
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Infof("deploy service stopping")
			return nil
		},
	})

	return svc
}

func HandlersProvider(
	logger *utils.Logger,
	ats services.AtsService,
	deploy services.DeployService,
	cfg *config.Config,
	redis *redis.RedisClient,
) *handlers.Handlers {
	return handlers.NewHandlers(ats, deploy, logger, cfg, redis)
}
