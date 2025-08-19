package service

import (
	"github.com/mohammadghasemi1379/sms-gateway/internal/entity"
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

func (s *userService) CreateUser(user *entity.User) error {

	return nil
}

func (s *userService) UpdateCredit(userID uint64, amount uint32) error {
	return nil
}
