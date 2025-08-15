package connection

import (
	"context"

	"github.com/mohammadghasemi1379/sms-gateway/config"
	"github.com/mohammadghasemi1379/sms-gateway/pkg/logger"

	"github.com/redis/go-redis/v9"
)

func RedisConnection(ctx context.Context, logger *logger.Logger, config config.Config) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Redis.GetRedisAddr(),
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
		PoolSize: 10000,
	})

	logger.Info(ctx, "Pinging main redis")

	redisStatus := redisClient.Ping(ctx)
	if redisStatus.Err() != nil {
		logger.Panic(ctx, "error on ping main redis", redisStatus.Err())
	}

	logger.Info(ctx, "Connected to main redis, status: %s", redisStatus.Val())

	return redisClient
}

func RedisPubSubConnection(ctx context.Context, logger *logger.Logger, config config.Config) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Redis.GetRedisAddr(),
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
		PoolSize: 10000,
	})

	logger.Info(ctx, "Pinging PUB SUB redis")

	redisStatus := redisClient.Ping(ctx)
	if redisStatus.Err() != nil {
		logger.Panic(ctx, "error on ping PUB SUB redis", redisStatus.Err())
	}

	logger.Info(ctx, "Connected to PUB SUB redis, status: %s", redisStatus.Val())

	return redisClient
}
