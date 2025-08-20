package errors

import (
	"errors"
	"strings"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

// Custom error types
var (
	ErrUserAlreadyExists = errors.New("user with this phone number already exists")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidInput      = errors.New("invalid input data")
)

// Business logic error types
type BusinessError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e BusinessError) Error() string {
	return e.Message
}

// Error codes
const (
	CodeUserAlreadyExists = "USER_ALREADY_EXISTS"
	CodeUserNotFound      = "USER_NOT_FOUND"
	CodeInvalidInput      = "INVALID_INPUT"
	CodeInternalError     = "INTERNAL_ERROR"
)

// NewBusinessError creates a new business error
func NewBusinessError(code, message string) *BusinessError {
	return &BusinessError{
		Code:    code,
		Message: message,
	}
}

// ParseDatabaseError parses database errors and returns appropriate business errors
func ParseDatabaseError(err error) error {
	if err == nil {
		return nil
	}

	// Handle GORM specific errors
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NewBusinessError(CodeUserNotFound, "User not found")
	}

	// Handle MySQL specific errors
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		switch mysqlErr.Number {
		case 1062: // Duplicate entry error
			return parseDuplicateKeyError(mysqlErr)
		case 1452: // Foreign key constraint fails
			return NewBusinessError(CodeInvalidInput, "Invalid reference data")
		}
	}

	// Return original error if not recognized
	return err
}

// parseDuplicateKeyError parses MySQL 1062 duplicate key error
func parseDuplicateKeyError(mysqlErr *mysql.MySQLError) error {
	message := mysqlErr.Message

	// Check if it's related to phone number uniqueness
	if strings.Contains(message, "phone_number") || strings.Contains(message, "idx_users_phone_number") {
		return NewBusinessError(CodeUserAlreadyExists, "A user with this phone number already exists")
	}

	// Generic duplicate key error
	return NewBusinessError(CodeUserAlreadyExists, "A user with this information already exists")
}

// IsBusinessError checks if an error is a business error
func IsBusinessError(err error) (*BusinessError, bool) {
	var businessErr *BusinessError
	if errors.As(err, &businessErr) {
		return businessErr, true
	}
	return nil, false
}
