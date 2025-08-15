package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SMSHandler struct{}

func NewSMSHandler() *SMSHandler {
	return &SMSHandler{}
}

func (h *SMSHandler) Send(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "hello",
	})
}
