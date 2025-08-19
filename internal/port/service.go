package port

import "github.com/mohammadghasemi1379/sms-gateway/internal/entity"

type SMSService interface {
	SendSMS(req *entity.SMS) error
	GetUserHistory(userID uint64) ([]entity.SMS, error)
}

type TransactionService interface {
	CreateTransaction(transaction *entity.Transaction) error
	GetTransactionBySMSID(smsID uint64) (*entity.Transaction, error)
	UpdateTransactionStatus(smsID uint64, status entity.TransactionStatusEnum) error
}

type UserService interface {
	CreateUser(user *entity.User) error
	UpdateCredit(userID uint64, amount uint32) error
}

type SMSConsumer interface {
	Consume(sms *entity.SMS) error
}