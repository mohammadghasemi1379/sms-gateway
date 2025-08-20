package service

import (
	"context"
	"errors"

	"github.com/mohammadghasemi1379/sms-gateway/connection"
	"github.com/mohammadghasemi1379/sms-gateway/internal/entity"
	"github.com/mohammadghasemi1379/sms-gateway/internal/port"
	"github.com/mohammadghasemi1379/sms-gateway/pkg/logger"
	"github.com/redis/go-redis/v9"
)

type smsService struct {
	smsRepo         port.SMSRepository
	userRepo        port.UserRepository
	transactionRepo port.TransactionRepository
	rabbitMQConnection *connection.RabbitMQConnection
	redisClient *redis.Client
	logger *logger.Logger
}

func NewSMSService(
	smsRepo port.SMSRepository,
	userRepo port.UserRepository,
	transactionRepo port.TransactionRepository,
	rabbitMQConnection *connection.RabbitMQConnection,
	redisClient *redis.Client,
	logger *logger.Logger,
) port.SMSService {
	return &smsService{
		smsRepo:         smsRepo,
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
		rabbitMQConnection: rabbitMQConnection,
		redisClient: redisClient,
		logger: logger,
	}
}

func (s *smsService) SendSMS(ctx context.Context, sms *entity.SMS) error {
	sms = s.CalculateCost(ctx, sms)

	hasEnoughCredit, err := s.userRepo.HasEnoughCredit(ctx, sms.UserID, sms.Cost)
	if err != nil {
		s.logger.Error(ctx, "failed to check if user has enough credit", "error", err)
		return err
	}

	if !hasEnoughCredit {
		return errors.New("user does not have enough credit")
	}

	err = s.smsRepo.Create(ctx, sms)
	if err != nil {
		s.logger.Error(ctx, "failed to create sms", "error", err)
		return err
	}

	transaction := &entity.Transaction{
		UserID: sms.UserID,
		Amount: float64(sms.Cost),
		Status: entity.TransactionPending,
		Operation: entity.Decrease,
		SMSID: &sms.ID,
	}
	err = s.transactionRepo.Create(ctx, transaction)
	if err != nil {
		s.logger.Error(ctx, "failed to create transaction", "error", err)
		return err
	}

	user, err := s.userRepo.GetByID(ctx, sms.UserID)
	if err != nil {
		s.logger.Error(ctx, "failed to get user", "error", err)
		return err
	}

	err = s.userRepo.DecreaseCredit(ctx, user, sms.Cost)
	if err != nil {
		s.logger.Error(ctx, "failed to decrease credit", "error", err)
		return err
	}

	return nil
}

func (s *smsService) GetUserHistory(ctx context.Context, userID uint64, page int, pageSize int) ([]entity.SMS, error) {
	return s.smsRepo.UserHistory(ctx, userID, pageSize, (page-1)*pageSize)
}

func (s *smsService) CalculateCost(ctx context.Context, sms *entity.SMS) *entity.SMS {
	// follow the project rules all the sms cost is fixed amount
	sms.Cost = 1000
	return sms
}
