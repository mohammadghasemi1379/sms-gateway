package service

import (
	"github.com/mohammadghasemi1379/sms-gateway/internal/entity"
	"github.com/mohammadghasemi1379/sms-gateway/internal/port"
)

type transactionService struct {
	transactionRepo port.TransactionRepository
	userRepo        port.UserRepository
}

func NewTransactionService(
	transactionRepo port.TransactionRepository,
	userRepo port.UserRepository,
) port.TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
	}
}

func (s *transactionService) CreateTransaction(transaction *entity.Transaction) error {
	return nil
}

func (s *transactionService) GetTransactionBySMSID(smsID uint64) (*entity.Transaction, error) {
	return nil, nil
}

func (s *transactionService) UpdateTransactionStatus(smsID uint64, status entity.TransactionStatusEnum) error {
	return nil
}
