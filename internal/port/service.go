package port

import (
	"context"

	"github.com/mohammadghasemi1379/sms-gateway/connection"
	"github.com/mohammadghasemi1379/sms-gateway/internal/entity"
)

type SMSService interface {
	SendSMS(ctx context.Context, sms *entity.SMS) error
	GetUserHistory(ctx context.Context, userID uint64, page int, pageSize int) ([]entity.SMS, error)
	CalculateCost(ctx context.Context, sms *entity.SMS) *entity.SMS
	GetSMSByID(ctx context.Context, smsID uint64) (*entity.SMS, error)
	UpdateSMSStatus(ctx context.Context, smsID uint64, status entity.SMSStatusEnum) error
}

type TransactionService interface {
	UpdateTransactionStatus(ctx context.Context, smsID uint64, status entity.TransactionStatusEnum) error
}

type UserService interface {
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	UpdateCredit(ctx context.Context, userID uint64, amount uint32) (*entity.User, error)
}

type MultiQueueConsumer interface {
	ConsumeAllQueues(ctx context.Context) error
}

type QueueManager interface {
	DetermineQueue(ctx context.Context) (string, error)
	PublishToQueue(ctx context.Context, message connection.RabbitMQMessageBody) error
	GetQueueNames() []string
}