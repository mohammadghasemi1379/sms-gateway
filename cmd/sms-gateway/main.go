package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mohammadghasemi1379/sms-gateway/config"
	"github.com/mohammadghasemi1379/sms-gateway/connection"
	"github.com/mohammadghasemi1379/sms-gateway/internal/handler"
	"github.com/mohammadghasemi1379/sms-gateway/internal/migration"
	"github.com/mohammadghasemi1379/sms-gateway/internal/repository"
	"github.com/mohammadghasemi1379/sms-gateway/internal/service"
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
	_ = connection.RedisPubSubConnection(ctx, logger, *cfg)
	gormDB, sqlDB := connection.MysqlConnection(ctx, logger, cfg)
	RabbitMQConnection := connection.NewRabbitMQConnection(cfg.RabbitMQ, logger, "sms-gateway", "sms-gateway", "sms-gateway")

	// Run database migrations
	migrationRunner := migration.NewRunner(gormDB, sqlDB, logger)
	if err := migrationRunner.RunMigrations(ctx, "migration"); err != nil {
		logger.Panic(ctx, "Failed to run database migrations", err)
	}

	smsRepository := repository.NewSMSRepository(gormDB, logger)
	transactionRepository := repository.NewTransactionRepository(gormDB, logger)
	userRepository := repository.NewUserRepository(gormDB, logger)

	smsService := service.NewSMSService(smsRepository, userRepository, transactionRepository, RabbitMQConnection, redisConnection, logger)
	userService := service.NewUserService(userRepository, transactionRepository)

	// Initialize Gin router
	if cfg.App.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	} else {
		gin.SetMode(gin.DebugMode)
	}
	router := gin.Default()

	// Initialize handlers
	smsHandler := handler.NewSMSHandler(smsService, logger)
	userHandler := handler.NewUserHandler(userService, logger)

	// Setup routes
	api := router.Group("/api")
	{
		user := api.Group("/user")
		{
			user.POST("/create", userHandler.CreateUser)
			user.POST("/update-credit", userHandler.UpdateCredit)
		}
		sms := api.Group("/sms")
		{
			sms.POST("/send", smsHandler.Send)
			sms.GET("/history", smsHandler.GetHistory)
		}
	}

	// Create HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.App.Host, cfg.App.Port),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		logger.Info(ctx, "Starting HTTP server", "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(ctx, "Failed to start server", "error", err.Error())
		}
	}()

	// Wait for shutdown signal
	<-ctx.Done()
	logger.Info(
		ctx,
		"Shutting down SMS Gateway service",
		"shutdown_time", time.Now(),
	)

	// Shutdown server gracefully
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error(ctx, "Server forced to shutdown", "error", err.Error())
	} else {
		logger.Info(ctx, "Server exited gracefully")
	}
}
