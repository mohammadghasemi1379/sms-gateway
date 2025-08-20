package service

import (
	"context"

	"github.com/mohammadghasemi1379/sms-gateway/internal/entity"
	"github.com/mohammadghasemi1379/sms-gateway/internal/errors"
	"github.com/mohammadghasemi1379/sms-gateway/internal/port"
)

type userService struct {
	userRepo        port.UserRepository
	transactionRepo port.TransactionRepository
}

func NewUserService(
	userRepo port.UserRepository,
	transactionRepo port.TransactionRepository,
) port.UserService {
	return &userService{
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *userService) CreateUser(ctx context.Context, user *entity.User) (*entity.User, error) {
	user, err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, errors.ParseDatabaseError(err)
	}

	return user, nil
}

func (s *userService) UpdateCredit(ctx context.Context, userID uint64, amount uint32) (*entity.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	err = s.userRepo.IncreaseCredit(ctx, user, amount)
	if err != nil {
		return nil, errors.ParseDatabaseError(err)
	}

	return user, nil
}
