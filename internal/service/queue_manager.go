package service

import (
	"context"
	"fmt"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/mohammadghasemi1379/sms-gateway/config"
	"github.com/mohammadghasemi1379/sms-gateway/connection"
	"github.com/mohammadghasemi1379/sms-gateway/pkg/logger"
)

const (
	QueueNameSMSMain      = "sms-gateway"           // Main queue (existing)
	QueueNameSMSPrimary   = "sms-gateway-primary"   // 90% of overflow traffic
	QueueNameSMSSecondary = "sms-gateway-secondary" // 10% of overflow traffic
)

type QueueDistributionStrategy struct {
	logger             *logger.Logger
	rabbitMQConnection *connection.RabbitMQConnection
	prefetchCount      int
	messageCounter     int64
	config             config.RabbitMQConfig
}

func NewQueueDistributionStrategy(
	logger *logger.Logger,
	rabbitMQConnection *connection.RabbitMQConnection,
	prefetchCount int,
	config config.RabbitMQConfig,
) *QueueDistributionStrategy {
	return &QueueDistributionStrategy{
		logger:             logger,
		rabbitMQConnection: rabbitMQConnection,
		prefetchCount:      prefetchCount,
		messageCounter:     0,
		config:             config,
	}
}

func (q *QueueDistributionStrategy) DetermineQueue(ctx context.Context) (string, error) {
	mainQueueCount, err := q.getQueueMessageCount(ctx, QueueNameSMSMain)
	if err != nil {
		q.logger.Error(ctx, "Failed to get main queue count", "error", err.Error())
		return QueueNameSMSMain, nil
	}

	q.logger.Debug(ctx, "Queue status",
		"main_queue_count", mainQueueCount,
		"prefetch_threshold", q.prefetchCount,
	)

	if mainQueueCount < int64(q.prefetchCount) {
		return QueueNameSMSMain, nil
	}

	return q.selectOverflowQueue(ctx), nil
}

func (q *QueueDistributionStrategy) selectOverflowQueue(ctx context.Context) string {
	rand.Seed(time.Now().UnixNano() + atomic.AddInt64(&q.messageCounter, 1))
	randomValue := rand.Intn(100) + 1

	if randomValue <= q.config.PrimaryWeight {
		q.logger.Debug(ctx, "Selected primary overflow queue", "random_value", randomValue)
		return QueueNameSMSPrimary
	}

	q.logger.Debug(ctx, "Selected secondary overflow queue", "random_value", randomValue)
	return QueueNameSMSSecondary
}

func (q *QueueDistributionStrategy) getQueueMessageCount(ctx context.Context, queueName string) (int64, error) {
	tempConn := connection.NewRabbitMQConnection(
		q.config,
		q.logger,
		"queue-inspector",
		"sms-gateway", 
		queueName,
	)

	if err := tempConn.Connect(); err != nil {
		return 0, fmt.Errorf("failed to connect to inspect queue: %w", err)
	}
	defer tempConn.Close()

	count, err := tempConn.GetQueueMessageCount()
	if err != nil {
		return 0, fmt.Errorf("failed to get queue message count: %w", err)
	}

	return int64(count), nil
}

func (q *QueueDistributionStrategy) PublishToQueue(ctx context.Context, message connection.RabbitMQMessageBody) error {
	targetQueue, err := q.DetermineQueue(ctx)
	if err != nil {
		return fmt.Errorf("failed to determine target queue: %w", err)
	}

	msg := connection.RabbitMQMessage{
		Queue:       targetQueue,
		ContentType: "text/plain",
		Body:        message,
	}

	q.logger.Info(ctx, "Publishing message to queue",
		"queue", targetQueue,
		"message_type", message.Type,
	)

	if err := q.rabbitMQConnection.Publish(ctx, msg); err != nil {
		return fmt.Errorf("failed to publish to queue %s: %w", targetQueue, err)
	}

	return nil
}

func (q *QueueDistributionStrategy) GetQueueNames() []string {
	return []string{
		QueueNameSMSMain,
		QueueNameSMSPrimary,
		QueueNameSMSSecondary,
	}
}
