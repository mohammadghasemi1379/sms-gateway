package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/mohammadghasemi1379/sms-gateway/internal/entity"
	"github.com/mohammadghasemi1379/sms-gateway/internal/port"
	"github.com/mohammadghasemi1379/sms-gateway/pkg/logger"
)

type SMSHandler struct {
	smsService port.SMSService
	logger *logger.Logger
}

func NewSMSHandler(smsService port.SMSService, logger *logger.Logger) *SMSHandler {
	return &SMSHandler{
		smsService: smsService,
		logger: logger,
	}
}


type SendSMSRequest struct {
	ReceiveNumber string `json:"phone_number" binding:"required"`
	Message       string `json:"message" binding:"required"`
	UserID        uint64 `json:"user_id" binding:"required"`
}

func (h *SMSHandler) Send(c *gin.Context) {
	var req SendSMSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sms := &entity.SMS{
		ReceiveNumber: req.ReceiveNumber,
		Message:       req.Message,
		UserID:        req.UserID,
		Status:        entity.SMSStatusPending,
	}

	err := h.smsService.SendSMS(c, sms)
	if err != nil {
		h.logger.Error(c, "failed to send sms", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "message in queue",
	})
}


type SMSHistoryRequest struct {
	UserID uint64 `json:"user_id" binding:"required"`
	Page int `json:"page"`
	PageSize int `json:"page_size"`
}

func (h *SMSHandler) GetHistory(c *gin.Context) {
	var req SMSHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}

	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	history, err := h.smsService.GetUserHistory(c, req.UserID, req.Page, req.PageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}
