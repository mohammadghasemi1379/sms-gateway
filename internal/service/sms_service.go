package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/mohammadghasemi1379/sms-gateway/connection"
	"github.com/mohammadghasemi1379/sms-gateway/internal/entity"
	"github.com/mohammadghasemi1379/sms-gateway/internal/port"
	"github.com/mohammadghasemi1379/sms-gateway/pkg/logger"
)

type smsService struct {
	smsRepo            port.SMSRepository
	userRepo           port.UserRepository
	transactionRepo    port.TransactionRepository
	rabbitMQConnection *connection.RabbitMQConnection
	logger             *logger.Logger
	queueStrategy      *QueueDistributionStrategy
}

func NewSMSService(
	smsRepo port.SMSRepository,
	userRepo port.UserRepository,
	transactionRepo port.TransactionRepository,
	rabbitMQConnection *connection.RabbitMQConnection,
	logger *logger.Logger,
	queueStrategy *QueueDistributionStrategy,
) port.SMSService {
	return &smsService{
		smsRepo:            smsRepo,
		userRepo:           userRepo,
		transactionRepo:    transactionRepo,
		rabbitMQConnection: rabbitMQConnection,
		logger:             logger,
		queueStrategy:      queueStrategy,
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
		UserID:    sms.UserID,
		Amount:    float64(sms.Cost),
		Status:    entity.TransactionPending,
		Operation: entity.Decrease,
		SMSID:     &sms.ID,
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

    err = s.queueStrategy.PublishToQueue(ctx, connection.RabbitMQMessageBody{
		Data: fmt.Appendf(nil, "%d", sms.ID),
		Type: "sms",
	})
    if err != nil {
        s.logger.Error(ctx, "failed to publish sms to queue", "error", err)
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

func (s *smsService) GetSMSByID(ctx context.Context, smsID uint64) (*entity.SMS, error) {
	return s.smsRepo.GetByID(ctx, smsID)
}

func (s *smsService) UpdateSMSStatus(ctx context.Context, smsID uint64, status entity.SMSStatusEnum) error {
	return s.smsRepo.UpdateStatus(ctx, smsID, status)
}