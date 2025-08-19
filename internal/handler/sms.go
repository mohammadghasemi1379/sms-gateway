package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SMSHandler struct{}

func NewSMSHandler() *SMSHandler {
	return &SMSHandler{}
}

type SendSMSRequest struct {
	ReceiveNumber string `json:"phone_number"`
	Message       string `json:"message"`
	UserID        string `json:"user_id"`
}

func (h *SMSHandler) Send(c *gin.Context) {
	var req SendSMSRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "hello",
	})
}
