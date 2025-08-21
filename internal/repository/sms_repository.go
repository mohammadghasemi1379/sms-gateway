package repository

import (
	"context"

	"github.com/mohammadghasemi1379/sms-gateway/internal/entity"
	"github.com/mohammadghasemi1379/sms-gateway/internal/port"
	"github.com/mohammadghasemi1379/sms-gateway/pkg/logger"
	"gorm.io/gorm"
)

type smsRepository struct {
	db *gorm.DB
	logger *logger.Logger
}

func NewSMSRepository(db *gorm.DB, logger *logger.Logger) port.SMSRepository {
	return &smsRepository{
		db: db,
		logger: logger,
	}
}

func (r *smsRepository) Create(ctx context.Context, sms *entity.SMS) error {
	err := r.db.WithContext(ctx).Create(sms).Error
	if err != nil {
		r.logger.Error(ctx, "Failed to create sms", "error", err.Error())
		return err
	}
	return nil
}

func (r *smsRepository) GetByID(ctx context.Context, id uint64) (*entity.SMS, error) {
	var sms entity.SMS
	err := r.db.WithContext(ctx).First(&sms, id).Error
	if err != nil {
		return nil, err
	}
	return &sms, nil
}

func (r *smsRepository) Update(ctx context.Context, sms *entity.SMS) error {
	err := r.db.WithContext(ctx).Save(sms).Error
	if err != nil {
		r.logger.Error(ctx, "Failed to update sms", "error", err.Error())
		return err
	}
	return nil
}

func (r *smsRepository) UserHistory(ctx context.Context, userID uint64, limit int, offset int) ([]entity.SMS, error) {
	var smsList []entity.SMS
	err := r.db.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Limit(limit).Offset(offset).Find(&smsList).Error
	if err != nil {
		r.logger.Error(ctx, "Failed to get user history", "error", err.Error())
		return nil, err
	}
	return smsList, err
}

func (r *smsRepository) UpdateStatus(ctx context.Context, smsID uint64, status entity.SMSStatusEnum) error {
	err := r.db.WithContext(ctx).Model(&entity.SMS{}).Where("id = ?", smsID).Update("status", status).Error
	if err != nil {
		r.logger.Error(ctx, "Failed to update sms status", "error", err.Error())
		return err
	}
	return nil
}