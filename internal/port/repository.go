package port

import (
	"context"

	"github.com/mohammadghasemi1379/sms-gateway/internal/entity"
)

type SMSRepository interface {
	Create(ctx context.Context, sms *entity.SMS) error
	GetByID(ctx context.Context, id uint64) (*entity.SMS, error)
	Update(ctx context.Context, sms *entity.SMS) error
	UserHistory(ctx context.Context, userID uint64) ([]entity.SMS, error)
}

type TransactionRepository interface {
	Create(ctx context.Context, transaction *entity.Transaction) error
	GetBySMSID(ctx context.Context, smsID uint64) (*entity.Transaction, error)
	UpdateStatusBySMSID(ctx context.Context, smsID uint64, status entity.TransactionStatusEnum) error
}

type UserRepository interface {
	GetByID(ctx context.Context, id uint64) (*entity.User, error)
	Create(ctx context.Context, user *entity.User) (*entity.User, error)
	HasEnoughCredit(ctx context.Context, userID uint64, amount uint32) (bool, error)
	IncreaseCredit(ctx context.Context, user *entity.User, amount uint32) error
	DecreaseCredit(ctx context.Context, user *entity.User, amount uint32) error
}
