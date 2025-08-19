package repository

import (
	"github.com/mohammadghasemi1379/sms-gateway/internal/entity"
	"github.com/mohammadghasemi1379/sms-gateway/internal/port"
	"gorm.io/gorm"
)

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) port.TransactionRepository {
	return &transactionRepository{
		db: db,
	}
}

func (r *transactionRepository) Create(transaction *entity.Transaction) error {
	return r.db.Create(transaction).Error
}

func (r *transactionRepository) GetBySMSID(smsID uint64) (*entity.Transaction, error) {
	var transaction entity.Transaction
	err := r.db.Where("sms_id = ?", smsID).First(&transaction).Error
	if err != nil {
		return nil, err
	}
	return &transaction, nil
}

func (r *transactionRepository) UpdateStatusBySMSID(smsID uint64, status entity.TransactionStatusEnum) error {
	return r.db.Model(&entity.Transaction{}).Where("sms_id = ?", smsID).Update("status", status).Error
}
