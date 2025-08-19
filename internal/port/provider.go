package port

import (
	"context"

	"github.com/mohammadghasemi1379/sms-gateway/internal/entity"
)


type Provider interface {
	Send(ctx context.Context, sms *entity.SMS) (any, error)
	DeliveryReport(ctx context.Context, sms *entity.SMS) (any, error)
}