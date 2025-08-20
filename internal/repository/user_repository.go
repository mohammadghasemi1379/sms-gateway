package repository

import (
	"context"

	"github.com/mohammadghasemi1379/sms-gateway/internal/entity"
	"github.com/mohammadghasemi1379/sms-gateway/internal/port"
	"github.com/mohammadghasemi1379/sms-gateway/pkg/logger"
	"gorm.io/gorm"
)

type userRepository struct {
	db     *gorm.DB
	logger *logger.Logger
}

func NewUserRepository(db *gorm.DB, logger *logger.Logger) port.UserRepository {
	return &userRepository{
		db:     db,
		logger: logger,
	}
}

func (r *userRepository) GetByID(ctx context.Context, id uint64) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).First(&user, id).Error
	if err != nil {
		r.logger.Error(ctx, "Failed to get user by id", "error", err.Error())
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) Create(ctx context.Context, user *entity.User) (*entity.User, error) {
	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		r.logger.Error(ctx, "Failed to create user", "error", err.Error())
		return nil, err
	}

	return user, nil
}

func (r *userRepository) HasEnoughCredit(ctx context.Context, userID uint64, amount uint32) (bool, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Select("credit").First(&user, userID).Error
	if err != nil {
		return false, err
	}
	return user.Credit >= int64(amount), nil
}

func (r *userRepository) IncreaseCredit(ctx context.Context, user *entity.User, amount uint32) error {
	err := r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", user.ID).Update("credit", gorm.Expr("credit + ?", amount)).Error
	if err != nil {
		r.logger.Error(ctx, "Failed to increase credit", "error", err.Error())
		return err
	}
	return nil
}

func (r *userRepository) DecreaseCredit(ctx context.Context, user *entity.User, amount uint32) error {
	err := r.db.WithContext(ctx).Model(&entity.User{}).Where("id = ?", user.ID).Update("credit", gorm.Expr("credit - ?", amount)).Error
	if err != nil {
		r.logger.Error(ctx, "Failed to decrease credit", "error", err.Error())
		return err
	}
	return nil
}
