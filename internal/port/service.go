package port

import (
	"context"
	"github.com/mohammadghasemi1379/sms-gateway/internal/entity"
)

type SMSService interface {
	SendSMS(ctx context.Context, sms *entity.SMS) error
	GetUserHistory(ctx context.Context, userID uint64, page int, pageSize int) ([]entity.SMS, error)
	CalculateCost(ctx context.Context, sms *entity.SMS) *entity.SMS
}

type TransactionService interface {
	CreateTransaction(transaction *entity.Transaction) error
	GetTransactionBySMSID(smsID uint64) (*entity.Transaction, error)
	UpdateTransactionStatus(smsID uint64, status entity.TransactionStatusEnum) error
}

type UserService interface {
	CreateUser(ctx context.Context, user *entity.User) (*entity.User, error)
	UpdateCredit(ctx context.Context, userID uint64, amount uint32) (*entity.User, error)
}

type SMSConsumer interface {
	Consume(ctx context.Context, sms *entity.SMS) error
}
