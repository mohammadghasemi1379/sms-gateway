package service

import (
	"context"

	"github.com/mohammadghasemi1379/sms-gateway/internal/entity"
	"github.com/mohammadghasemi1379/sms-gateway/internal/port"
	"github.com/mohammadghasemi1379/sms-gateway/pkg/logger"
)

type transactionService struct {
	transactionRepo port.TransactionRepository
	userRepo        port.UserRepository
	logger *logger.Logger
}

func NewTransactionService(
	transactionRepo port.TransactionRepository,
	userRepo port.UserRepository,
	logger *logger.Logger,
) port.TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
		logger: logger,
	}
}


func (s *transactionService) UpdateTransactionStatus(ctx context.Context, smsID uint64, status entity.TransactionStatusEnum) error {
	return s.transactionRepo.UpdateStatusBySMSID(ctx, smsID, status)
}
