package entity

import "time"

type OperationEnum string

const (
	Increase OperationEnum = "INCREASE"
	Decrease OperationEnum = "DECREASE"
)

type TransactionStatusEnum string

const (
	TransactionPending TransactionStatusEnum = "PENDING"
	TransactionFailed  TransactionStatusEnum = "FAILED"
	TransactionSuccess TransactionStatusEnum = "SUCCESS"
)

type Transaction struct {
	ID        uint64                `json:"id"`
	UserID    uint64                `json:"user_id"`
	SMSID     *uint64               `json:"sms_id"` // Optional, for SMS-related transactions
	Amount    float64               `json:"amount"`
	Operation OperationEnum         `json:"operation"`
	Status    TransactionStatusEnum `json:"status"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}
