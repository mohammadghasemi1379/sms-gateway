package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/mohammadghasemi1379/sms-gateway/config"
	"github.com/mohammadghasemi1379/sms-gateway/connection"
	"github.com/mohammadghasemi1379/sms-gateway/pkg/logger"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer cancel()

	logger := logger.New()
	cfg := config.Load()

	logger.Info(ctx, "Starting SMS Gateway service",
		"service", cfg.App.Name,
		"version", cfg.App.Version,
		"environment", cfg.App.Environment,
		"port", cfg.App.Port,
		"start_time", time.Now(),
	)

	redisConnection := connection.RedisConnection(ctx, logger, *cfg)
	redisPubSubConnection := connection.RedisPubSubConnection(ctx, logger, *cfg)
	mysqlConnection, _ := connection.MysqlConnection(ctx, logger, cfg)

	// Wait for shutdown signal
	<-ctx.Done()
	logger.Info(
		ctx,
		"Shutting down SMS Gateway service",
		"shutdown_time", time.Now(),
	)
}
