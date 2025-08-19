package service

import (
	"github.com/mohammadghasemi1379/sms-gateway/internal/entity"
	"github.com/mohammadghasemi1379/sms-gateway/internal/port"
)

type smsConsumer struct {
	smsRepo         port.SMSRepository
	transactionRepo port.TransactionRepository
}

func NewSMSConsumer(
	smsRepo port.SMSRepository,
	transactionRepo port.TransactionRepository,
) port.SMSConsumer {
	return &smsConsumer{
		smsRepo:         smsRepo,
		transactionRepo: transactionRepo,
	}
}

func (c *smsConsumer) Consume(sms *entity.SMS) error {
	return nil
}
