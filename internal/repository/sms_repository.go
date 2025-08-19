package repository

import (
	"github.com/mohammadghasemi1379/sms-gateway/internal/entity"
	"github.com/mohammadghasemi1379/sms-gateway/internal/port"
	"gorm.io/gorm"
)

type smsRepository struct {
	db *gorm.DB
}

func NewSMSRepository(db *gorm.DB) port.SMSRepository {
	return &smsRepository{
		db: db,
	}
}

func (r *smsRepository) Create(sms *entity.SMS) error {
	return r.db.Create(sms).Error
}

func (r *smsRepository) GetByID(id uint64) (*entity.SMS, error) {
	var sms entity.SMS
	err := r.db.First(&sms, id).Error
	if err != nil {
		return nil, err
	}
	return &sms, nil
}

func (r *smsRepository) Update(sms *entity.SMS) error {
	return r.db.Save(sms).Error
}

func (r *smsRepository) UserHistory(userID uint64) ([]entity.SMS, error) {
	var smsList []entity.SMS
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&smsList).Error
	return smsList, err
}
