package service

import (
	"errors"
	"github.com/mohammadghasemi1379/sms-gateway/internal/entity"
	"github.com/mohammadghasemi1379/sms-gateway/internal/port"
)

type smsService struct {
	smsRepo         port.SMSRepository
	userRepo        port.UserRepository
	transactionRepo port.TransactionRepository
}

func NewSMSService(
	smsRepo port.SMSRepository,
	userRepo port.UserRepository,
	transactionRepo port.TransactionRepository,
) port.SMSService {
	return &smsService{
		smsRepo:         smsRepo,
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *smsService) SendSMS(sms *entity.SMS) error {
	sms = s.CalculateCost(sms)	
	hasEnoughCredit, err := s.userRepo.HasEnoughCredit(sms.UserID, sms.Cost)
	if err != nil {
		return err
	}

	if !hasEnoughCredit {
		return errors.New("user does not have enough credit")
	}

	transaction := &entity.Transaction{
		UserID: sms.UserID,
		Amount: float64(sms.Cost),
		Status: entity.TransactionPending,
		Operation: entity.Decrease,
		SMSID: &sms.ID,
	}
	err = s.transactionRepo.Create(transaction)
	if err != nil {
		return err
	}

	err = s.userRepo.DecreaseCredit(sms.UserID, sms.Cost)
	if err != nil {
		return err
	}

	return nil
}

func (s *smsService) GetUserHistory(userID uint64) ([]entity.SMS, error) {
	return s.smsRepo.UserHistory(userID)
}

func (s *smsService) CalculateCost(sms *entity.SMS) *entity.SMS {
	// follow the project rules all the sms cost is fixed amount
	sms.Cost = 1000
	return sms
}
