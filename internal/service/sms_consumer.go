package service

import (
	"context"
	"fmt"
	"strconv"

	"github.com/mohammadghasemi1379/sms-gateway/connection"
	"github.com/mohammadghasemi1379/sms-gateway/internal/port"
	"github.com/mohammadghasemi1379/sms-gateway/pkg/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type smsConsumer struct {
	smsService         port.SMSService
	TransactionService port.TransactionService
	userService port.UserService
	rabbitMQConnection *connection.RabbitMQConnection
	provider        port.Provider
	logger          *logger.Logger
	prefetchCount int
}

func NewSMSConsumer(
	smsService port.SMSService,
	TransactionService port.TransactionService,
	userService port.UserService,
	rabbitMQConnection *connection.RabbitMQConnection,
	provider port.Provider,
	logger *logger.Logger,
	prefetchCount int,
) port.SMSConsumer {
	return &smsConsumer{
		smsService: smsService,
		TransactionService: TransactionService,
		userService: userService,
		rabbitMQConnection: rabbitMQConnection,
		provider: provider,
		logger: logger,
		prefetchCount: prefetchCount,
	}
}

func (c *smsConsumer) Consume(ctx context.Context) error {
	c.rabbitMQConnection.HandleConsumedDeliveries(ctx, c.prefetchCount, func(ctx context.Context, logger *logger.Logger, conn connection.RabbitMQConnection, delivery amqp.Delivery) {
		
		smsID, err := strconv.ParseUint(string(delivery.Body), 10, 64)
		if err != nil {
			c.logger.Error(ctx, "failed to parse sms id", "error", err.Error())
			//reject the message
			if nackErr := delivery.Nack(false, false); nackErr != nil {
				c.logger.Error(ctx, "Failed to Nack message after unmarshalling error", nackErr.Error())
			}
		}else{
			if err := c.processMessage(ctx, smsID); err != nil {
				c.logger.Error(ctx, "Failed to process sms message", err.Error())
				if nackErr := delivery.Nack(false, true); nackErr != nil {
					c.logger.Error(ctx, "failed to Nack", nackErr.Error())
				}
			}
			if AckErr := delivery.Ack(false); AckErr != nil {
				c.logger.Error(ctx, "Error on ack message", AckErr.Error())
			}
			c.logger.Info(ctx, fmt.Sprintf("Received delivery of sms %d", delivery.DeliveryTag))
		}
	})

	c.logger.Info(ctx, "[*] Waiting for messages. To exit press CTRL+C")
	return nil
}

func(c *smsConsumer) processMessage(ctx context.Context, smsID uint64) error {
	sms, err := c.smsService.GetSMSByID(ctx, smsID)
	if err != nil {
		c.logger.Error(ctx, "failed to get sms", "error", err.Error())
		return err
	}

	response, err := c.provider.Send(ctx, sms)
	if err != nil {
		c.logger.Error(ctx, "failed to send sms", "error", err.Error())
		return err
	}

	if response.Status == "ok" && response.Message == "sended" {
		return nil
	}

	c.logger.Error(ctx, "failed to send sms", "error", response.Message)
	return fmt.Errorf("failed to send sms: %s", response.Message)
}
