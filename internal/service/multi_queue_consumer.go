package service

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/mohammadghasemi1379/sms-gateway/config"
	"github.com/mohammadghasemi1379/sms-gateway/connection"
	"github.com/mohammadghasemi1379/sms-gateway/internal/entity"
	"github.com/mohammadghasemi1379/sms-gateway/internal/port"
	"github.com/mohammadghasemi1379/sms-gateway/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MultiQueueConsumer struct {
	smsService         port.SMSService
	transactionService port.TransactionService
	userService        port.UserService
	rabbitMQConnection *connection.RabbitMQConnection
	provider           port.Provider
	logger             *logger.Logger
	queueStrategy      *QueueDistributionStrategy
	config             config.RabbitMQConfig
}

func NewMultiQueueConsumer(
	smsService port.SMSService,
	transactionService port.TransactionService,
	userService port.UserService,
	rabbitMQConnection *connection.RabbitMQConnection,
	provider port.Provider,
	logger *logger.Logger,
	prefetchCount int,
	config config.RabbitMQConfig,
) *MultiQueueConsumer {
	return &MultiQueueConsumer{
		smsService:         smsService,
		transactionService: transactionService,
		userService:        userService,
		rabbitMQConnection: rabbitMQConnection,
		provider:           provider,
		logger:             logger,
		queueStrategy:      NewQueueDistributionStrategy(logger, rabbitMQConnection, prefetchCount, config),
		config:             config,
	}
}

func (c *MultiQueueConsumer) ConsumeAllQueues(ctx context.Context) error {
	if err := c.initializeQueues(); err != nil {
		return fmt.Errorf("failed to initialize queues: %w", err)
	}

	var wg sync.WaitGroup
	queueNames := c.queueStrategy.GetQueueNames()

	for _, queueName := range queueNames {
		wg.Add(1)
		go func(queue string) {
			defer wg.Done()
			c.logger.Info(ctx, "Starting consumer for queue", "queue", queue)

			if err := c.consumeQueue(ctx, queue); err != nil {
				c.logger.Error(ctx, "Failed to consume queue", "queue", queue, "error", err.Error())
			}
		}(queueName)
	}

	wg.Wait()
	return nil
}

func (c *MultiQueueConsumer) initializeQueues() error {
	queueNames := c.queueStrategy.GetQueueNames()
	return c.rabbitMQConnection.DeclareMultipleQueues(queueNames)
}

func (c *MultiQueueConsumer) consumeQueue(ctx context.Context, queueName string) error {
	queueConn := connection.NewRabbitMQConnection(
		c.config,
		c.logger,
		fmt.Sprintf("consumer-%s", queueName),
		"sms-gateway",
		queueName,
	)

	if err := queueConn.Connect(); err != nil {
		return fmt.Errorf("failed to connect to queue %s: %w", queueName, err)
	}
	defer queueConn.Close()

	queueConn.ConnectionOpener()

	queueConn.HandleConsumedDeliveries(ctx, c.config.PrefetchCount, func(ctx context.Context, logger *logger.Logger, conn connection.RabbitMQConnection, delivery amqp.Delivery) {
		c.processMessage(ctx, delivery, queueName)
	})

	return nil
}

func (c *MultiQueueConsumer) processMessage(ctx context.Context, delivery amqp.Delivery, queueName string) {
	c.logger.Info(ctx, "Processing message", "queue", queueName, "delivery_tag", delivery.DeliveryTag)

	smsID, err := strconv.ParseUint(string(delivery.Body), 10, 64)
	if err != nil {
		c.logger.Error(ctx, "failed to parse sms id", "error", err.Error(), "queue", queueName)
		if nackErr := delivery.Nack(false, false); nackErr != nil {
			c.logger.Error(ctx, "Failed to Nack message after unmarshalling error", nackErr.Error())
		}
		return
	}

	if err := c.processMessageLogic(ctx, smsID); err != nil {
		c.logger.Error(ctx, "Failed to process sms message", "error", err.Error(), "queue", queueName)
		if nackErr := delivery.Nack(false, true); nackErr != nil {
			c.logger.Error(ctx, "failed to Nack", nackErr.Error())
		}
		return
	}

	if ackErr := delivery.Ack(false); ackErr != nil {
		c.logger.Error(ctx, "Error on ack message", ackErr.Error())
	}

	c.logger.Info(ctx, "Message processed successfully", "queue", queueName, "sms_id", smsID)
}

func (c *MultiQueueConsumer) processMessageLogic(ctx context.Context, smsID uint64) error {
	sms, err := c.smsService.GetSMSByID(ctx, smsID)
	if err != nil {
		c.logger.Error(ctx, "failed to get sms", "error", err.Error())
		return err
	}

	response, err := c.provider.Send(ctx, sms)
	if err != nil {
		c.logger.Error(ctx, "provider failed to respond", "error", err.Error(), "sms_id", smsID)
		return err
	}

	if response.Status == "ok" && response.Message == "sended" {
		c.smsService.UpdateSMSStatus(ctx, smsID, entity.SMSStatusSent)
		c.transactionService.UpdateTransactionStatus(ctx, smsID, entity.TransactionSuccess)
		return nil
	}

	c.logger.Error(ctx, "provider returned not ok response", "status", response.Status, "message", response.Message, "sms_id", smsID)
	return fmt.Errorf("provider returned not ok response: status=%s, message=%s", response.Status, response.Message)
}
