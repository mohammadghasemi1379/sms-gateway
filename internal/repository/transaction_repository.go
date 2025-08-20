package repository

import (
	"context"

	"github.com/mohammadghasemi1379/sms-gateway/internal/entity"
	"github.com/mohammadghasemi1379/sms-gateway/internal/port"
	"github.com/mohammadghasemi1379/sms-gateway/pkg/logger"
	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
	logger *logger.Logger
}

func NewTransactionRepository(db *gorm.DB, logger *logger.Logger) port.TransactionRepository {
	return &transactionRepository{
		db: db,
		logger: logger,
	}
}

func (r *transactionRepository) Create(ctx context.Context, transaction *entity.Transaction) error {
	err := r.db.WithContext(ctx).Create(transaction).Error
	if err != nil {
		r.logger.Error(ctx, "Failed to create transaction", "error", err.Error())
		return err
	}
	return nil
}

func (r *transactionRepository) GetBySMSID(ctx context.Context, smsID uint64) (*entity.Transaction, error) {
	var transaction entity.Transaction
	err := r.db.WithContext(ctx).Where("sms_id = ?", smsID).First(&transaction).Error
	if err != nil {
		r.logger.Error(ctx, "Failed to get transaction by sms id", "error", err.Error())
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) UpdateStatusBySMSID(ctx context.Context, smsID uint64, status entity.TransactionStatusEnum) error {
	err := r.db.WithContext(ctx).Model(&entity.Transaction{}).Where("sms_id = ?", smsID).Update("status", status).Error
	if err != nil {
		r.logger.Error(ctx, "Failed to update transaction status by sms id", "error", err.Error())
		return err
	}
	return nil
}
