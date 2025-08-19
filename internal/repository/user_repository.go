package repository

import (
	"github.com/mohammadghasemi1379/sms-gateway/internal/entity"
	"github.com/mohammadghasemi1379/sms-gateway/internal/port"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) port.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) GetByID(id uint64) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(user *entity.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) HasEnoughCredit(userID uint64, amount uint32) (bool, error) {
	var user entity.User
	err := r.db.Select("credit").First(&user, userID).Error
	if err != nil {
		return false, err
	}
	return user.Credit >= int64(amount), nil
}

func (r *userRepository) UpdateCredit(userID uint64, amount uint32) error {
	return r.db.Model(&entity.User{}).Where("id = ?", userID).Update("credit", gorm.Expr("credit - ?", amount)).Error
}
