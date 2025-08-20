package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/mohammadghasemi1379/sms-gateway/internal/entity"
	"github.com/mohammadghasemi1379/sms-gateway/internal/port"
)

type SMSHandler struct {
	smsService port.SMSService
}

func NewSMSHandler(smsService port.SMSService) *SMSHandler {
	return &SMSHandler{
		smsService: smsService,
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "message in queue",
	})
}


type SMSHistoryRequest struct {
	UserID uint64 `json:"user_id"`
}

func (h *SMSHandler) GetHistory(c *gin.Context) {
	var req SMSHistoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	history, err := h.smsService.GetUserHistory(c, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}
