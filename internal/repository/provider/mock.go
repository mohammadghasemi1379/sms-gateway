package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/mohammadghasemi1379/sms-gateway/config"
	"github.com/mohammadghasemi1379/sms-gateway/internal/entity"
	"github.com/mohammadghasemi1379/sms-gateway/internal/port"
	"github.com/mohammadghasemi1379/sms-gateway/pkg/logger"
)

type ProviderMock struct {
	logger *logger.Logger
	config *config.Config
	client *http.Client
}

func NewMockProvider(logger *logger.Logger, config *config.Config) port.Provider {
	return &ProviderMock{
		logger: logger,
		config: config,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

type SendResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func (p *ProviderMock) Send(ctx context.Context, sms *entity.SMS) (any, error) {
	p.logger.Info(ctx, "Sending SMS", "sms", sms)

	url := fmt.Sprintf("%s:%d/mock/sms", p.config.Mock.Host, p.config.Mock.Port)

	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}

	res, err := p.client.Do(req.WithContext(ctx))
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		BodyErr := Body.Close()
		if BodyErr != nil {
			p.logger.Panic(ctx, "error in closing response body", err)
		}
	}(res.Body)

	if res.StatusCode != http.StatusOK {
		p.logger.Error(ctx, "error in sending sms", err)
		return nil, err
	}

	var response SendResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		p.logger.Error(ctx, "error in decode get credit", err)
		return nil, err
	}

	return response, nil
}

func (p *ProviderMock) DeliveryReport(ctx context.Context, sms *entity.SMS) (any, error) {
	return nil, nil
}
