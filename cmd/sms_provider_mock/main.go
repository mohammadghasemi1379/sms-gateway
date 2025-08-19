package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mohammadghasemi1379/sms-gateway/config"
	"github.com/mohammadghasemi1379/sms-gateway/pkg/logger"
)

var (
	requests = 0
	mu       sync.Mutex
	limit    = 5
	reset    = time.Now().Add(time.Minute)
)

type RequestBody struct {}

func handler(c *gin.Context) {
	mu.Lock()
	defer mu.Unlock()

	var req RequestBody
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if time.Now().After(reset) {
		requests = 0
		reset = time.Now().Add(time.Minute)
	}

	if requests >= limit {
		c.JSON(http.StatusTooManyRequests, gin.H{
			"status":  "fail",
			"message": "Too many requests",
		})
		return
	}
	requests++

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "sended",
	})
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer cancel()
	logger := logger.New()
	cfg := config.Load()

	// Initialize Gin router
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Setup routes
	api := router.Group("/mock")
	{
		api.POST("/sms", handler)
	}

	// Create HTTP server
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", cfg.Mock.Host, cfg.Mock.Port),
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		logger.Info(ctx, "Starting HTTP server", "address", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error(ctx, "Failed to start server", "error", err)
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
		logger.Error(ctx, "Server forced to shutdown", "error", err)
	} else {
		logger.Info(ctx, "Server exited gracefully")
	}
}
