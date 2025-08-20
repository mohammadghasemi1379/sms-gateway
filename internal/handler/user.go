package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mohammadghasemi1379/sms-gateway/internal/entity"
	"github.com/mohammadghasemi1379/sms-gateway/internal/errors"
	"github.com/mohammadghasemi1379/sms-gateway/internal/port"
	"github.com/mohammadghasemi1379/sms-gateway/pkg/logger"
)

type UserHandler struct {
	userService port.UserService
	logger *logger.Logger
}

func NewUserHandler(userService port.UserService, logger *logger.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger: logger,
	}
}

type CreateUserRequest struct {
	Name        string `json:"name" binding:"required,max=255"`
	PhoneNumber string `json:"phone_number" binding:"required,min=11,max=11"`
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user := &entity.User{
		Name:        req.Name,
		PhoneNumber: req.PhoneNumber,
	}

	user, err := h.userService.CreateUser(c, user)
	if err != nil {
		if businessErr, isBusiness := errors.IsBusinessError(err); isBusiness {
			statusCode := getHTTPStatusFromErrorCode(businessErr.Code)
			c.JSON(statusCode, gin.H{
				"error": businessErr.Message,
				"code":  businessErr.Code,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
			"code":  errors.CodeInternalError,
		})
		h.logger.Error(c, "Failed to create user", "error", err.Error())
		return
	}

	c.JSON(http.StatusCreated, user)
}

type UpdateCreditRequest struct {
	UserID uint64 `json:"user_id" binding:"required"`
	Amount uint32 `json:"amount" binding:"required"`
}


func (h *UserHandler) UpdateCredit(c *gin.Context) {
	var req UpdateCreditRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.userService.UpdateCredit(c, req.UserID, req.Amount)
	if err != nil {
		// Check if it's a business error
		if businessErr, isBusiness := errors.IsBusinessError(err); isBusiness {
			statusCode := getHTTPStatusFromErrorCode(businessErr.Code)
			c.JSON(statusCode, gin.H{
				"error": businessErr.Message,
				"code":  businessErr.Code,
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// getHTTPStatusFromErrorCode maps business error codes to HTTP status codes
func getHTTPStatusFromErrorCode(code string) int {
	switch code {
	case errors.CodeUserAlreadyExists:
		return http.StatusConflict
	case errors.CodeUserNotFound:
		return http.StatusNotFound
	case errors.CodeInvalidInput:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
