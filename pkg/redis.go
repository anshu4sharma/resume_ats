package redis

import (
	"context"
	"errors"
	"time"

	"github.com/anshu4sharma/resume_ats/pkg/utils"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

type RedisClient struct {
	*redis.Client
	logger  *utils.Logger
	isReady bool
}

// NewRedisClient creates a new RedisClient
func NewRedisClient(url string, logger *utils.Logger) *RedisClient {
	opt, err := redis.ParseURL(url)
	if err != nil {
		logger.Errorf(
			"invalid redis url",
			zap.String("url", url),
			zap.Error(err),
		)
		return nil
	}

	return &RedisClient{
		Client: redis.NewClient(opt),
		logger: logger,
	}
}

// Connect establishes connection with retry logic
func (r *RedisClient) Connect(maxRetries int, retryDelay time.Duration) error {
	var err error

	for i := 0; i < maxRetries; i++ {
		_, err = r.Ping(context.Background()).Result()
		if err == nil {
			r.isReady = true
			r.logger.Infof(
				"connected to redis",
				zap.Int("attempt", i+1),
			)
			return nil
		}

		r.logger.Errorf(
			"failed to connect to redis",
			zap.Int("attempt", i+1),
			zap.Int("max_retries", maxRetries),
			zap.Duration("retry_delay", retryDelay),
			zap.Error(err),
		)

		time.Sleep(retryDelay)
		retryDelay *= 2
	}

	return errors.New("max retries reached, could not connect to redis")
}

// Close closes the Redis connection
func (r *RedisClient) Close() error {
	if r.Client == nil {
		return nil
	}

	if err := r.Client.Close(); err != nil {
		r.logger.Errorf(
			"failed to close redis connection",
			zap.Error(err),
		)
		return err
	}

	r.logger.Infof("redis connection closed")
	return nil
}

// IsReady checks if Redis is connected
func (r *RedisClient) IsReady() bool {
	return r.isReady
}

// Wrapper methods

func (r *RedisClient) GetValue(ctx context.Context, key string) (string, error) {
	if !r.isReady {
		return "", errors.New("redis is not connected")
	}

	val, err := r.Get(ctx, key).Result()
	if err != nil {
		r.logger.Errorf(
			"failed to get redis key",
			zap.String("key", key),
			zap.Error(err),
		)
		return "", err
	}

	return val, nil
}

func (r *RedisClient) SetValue(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if !r.isReady {
		return errors.New("redis is not connected")
	}

	if err := r.Set(ctx, key, value, expiration).Err(); err != nil {
		r.logger.Errorf(
			"failed to set redis key",
			zap.String("key", key),
			zap.Duration("expiration", expiration),
			zap.Error(err),
		)
		return err
	}

	return nil
}

func (r *RedisClient) DeleteKey(ctx context.Context, key string) error {
	if !r.isReady {
		return errors.New("redis is not connected")
	}

	if err := r.Del(ctx, key).Err(); err != nil {
		r.logger.Errorf(
			"failed to delete redis key",
			zap.String("key", key),
			zap.Error(err),
		)
		return err
	}

	return nil
}
