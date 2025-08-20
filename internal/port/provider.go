package port

import (
	"context"

	"github.com/mohammadghasemi1379/sms-gateway/internal/entity"
)

type SendResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type Provider interface {
	Send(ctx context.Context, sms *entity.SMS) (*SendResponse, error)
	DeliveryReport(ctx context.Context, sms *entity.SMS) (any, error)
}