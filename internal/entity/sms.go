package entity

import "time"

type SMSStatusEnum string

const (
	SMSStatusPending SMSStatusEnum = "PENDING"
	SMSStatusSent    SMSStatusEnum = "SENT"
	SMSStatusFailed  SMSStatusEnum = "FAILED"
)

type SMS struct {
	ID            uint64        `json:"id"`
	UserID        uint64        `json:"user_id"`
	ReceiveNumber string        `json:"receive_number"`
	Message       string        `json:"message"`
	Status        SMSStatusEnum `json:"status"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}
