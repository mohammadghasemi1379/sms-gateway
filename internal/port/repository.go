package port

import (
	"github.com/mohammadghasemi1379/sms-gateway/internal/entity"
)

type SMSRepository interface {
	Create(sms *entity.SMS) error
	GetByID(id uint64) (*entity.SMS, error)
	Update(sms *entity.SMS) error
	UserHistory(userID uint64) ([]entity.SMS, error)
}

type TransactionRepository interface {
	Create(transaction *entity.Transaction) error
	GetBySMSID(smsID uint64) (*entity.Transaction, error)
	UpdateStatusBySMSID(smsID uint64, status entity.TransactionStatusEnum) error
}

type UserRepository interface {
	GetByID(id uint64) (*entity.User, error)
	Create(user *entity.User) error
	HasEnoughCredit(userID uint64, amount uint32) (bool, error)
	IncreaseCredit(userID uint64, amount uint32) error
	DecreaseCredit(userID uint64, amount uint32) error
}
