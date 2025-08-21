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
	"github.com/mohammadghasemi1379/sms-gateway/internal/repository/provider"
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

	gormDB, sqlDB := connection.MysqlConnection(ctx, logger, cfg)

	RabbitMQConnection := connection.NewRabbitMQConnection(cfg.RabbitMQ, logger, "sms-gateway", "sms-gateway", "sms-gateway")
	err := RabbitMQConnection.Connect()
	if err != nil {
		logger.Error(ctx, "Failed to connect to RabbitMQ", err)
	}
	RabbitMQConnection.ConnectionOpener()
	defer func(rabbitmqConn *connection.RabbitMQConnection) {
		err := rabbitmqConn.Close()
		if err != nil {
			logger.Error(ctx, "error on close RabbitMQ connection", err.Error())
		}
	}(RabbitMQConnection)

	queueStrategy := service.NewQueueDistributionStrategy(logger, RabbitMQConnection, cfg.RabbitMQ.PrefetchCount, cfg.RabbitMQ)
	if err := RabbitMQConnection.DeclareMultipleQueues(queueStrategy.GetQueueNames()); err != nil {
		logger.Panic(ctx, "Failed to declare queues", err)
	}

	// Run database migrations
	migrationRunner := migration.NewRunner(gormDB, sqlDB, logger)
	if err := migrationRunner.RunMigrations(ctx, "migration"); err != nil {
		logger.Panic(ctx, "Failed to run database migrations", err)
	}

	// Initialize repositories
	smsRepository := repository.NewSMSRepository(gormDB, logger)
	transactionRepository := repository.NewTransactionRepository(gormDB, logger)
	userRepository := repository.NewUserRepository(gormDB, logger)

	// Initialize services
	smsService := service.NewSMSService(smsRepository, userRepository, transactionRepository, RabbitMQConnection, logger, queueStrategy)
	userService := service.NewUserService(userRepository, transactionRepository)
	transactionService := service.NewTransactionService(transactionRepository, userRepository, logger)
	provider := provider.NewMockProvider(logger, cfg)

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

	// Start multi-queue consumer
	multiQueueConsumer := service.NewMultiQueueConsumer(smsService, transactionService, userService, RabbitMQConnection, provider, logger, cfg.RabbitMQ.PrefetchCount, cfg.RabbitMQ)
    go func() {
        logger.Info(ctx, "Starting multi-queue SMS consumer...")
        if err := multiQueueConsumer.ConsumeAllQueues(ctx); err != nil {
            logger.Error(ctx, "Failed to start multi-queue consumer", "error", err.Error())
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
		logger.Error(shutdownCtx, "Server forced to shutdown", "error", err.Error())
	} else {
		logger.Info(shutdownCtx, "Server exited gracefully")
	}
}
