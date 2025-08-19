package service

import (
	"context"
	"fmt"

	"github.com/mohammadghasemi1379/sms-gateway/internal/entity"
	"github.com/mohammadghasemi1379/sms-gateway/internal/port"
	"github.com/mohammadghasemi1379/sms-gateway/internal/repository/provider"
	"github.com/mohammadghasemi1379/sms-gateway/pkg/logger"
)

type smsConsumer struct {
	smsRepo         port.SMSRepository
	transactionRepo port.TransactionRepository
	provider        port.Provider
	logger          *logger.Logger
}

func NewSMSConsumer(
	smsRepo port.SMSRepository,
	transactionRepo port.TransactionRepository,
	provider port.Provider,
) port.SMSConsumer {
	return &smsConsumer{
		smsRepo:         smsRepo,
		transactionRepo: transactionRepo,
		provider:        provider,
	}
}

func (c *smsConsumer) Consume(ctx context.Context, sms *entity.SMS) error {
	response, err := c.provider.Send(ctx, sms)
	if err != nil {
		c.logger.Error(ctx, "failed to send sms", "error", err)
		return err
	}

	fmt.Println(response)

	sendResponse, ok := response.(*provider.SendResponse)
	if !ok {
		c.logger.Error(ctx, "invalid response type", "error", response)
		return fmt.Errorf("invalid response type")
	}

	if sendResponse.Status == "success" {
		return nil
	}

	c.logger.Error(ctx, "failed to send sms", "error", sendResponse.Message)
	return fmt.Errorf("failed to send sms: %s", sendResponse.Message)
}
