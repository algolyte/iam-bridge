package middleware

import (
	"errors"
	"github.com/zahidhasanpapon/iam-bridge/internal/provider"
	"net/http"

	"github.com/gin-gonic/gin"
)

// APIError represents a standardized API error response
type APIError struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

// ErrorHandlerMiddleware handles errors in a standardized way
func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		// Check if there are any errors
		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err
			handleError(c, err)
		}
	}
}

// handleError processes different types of errors and returns appropriate responses
func handleError(c *gin.Context, err error) {
	requestID := GetRequestID(c)

	// Map common errors to HTTP status codes and error codes
	switch {
	case errors.Is(err, provider.ErrInvalidCredentials):
		c.JSON(http.StatusUnauthorized, APIError{
			Code:      "INVALID_CREDENTIALS",
			Message:   "Invalid username or password",
			RequestID: requestID,
		})

	case errors.Is(err, provider.ErrTokenExpired):
		c.JSON(http.StatusUnauthorized, APIError{
			Code:      "TOKEN_EXPIRED",
			Message:   "Authentication token has expired",
			RequestID: requestID,
		})

	case errors.Is(err, provider.ErrTokenInvalid):
		c.JSON(http.StatusUnauthorized, APIError{
			Code:      "INVALID_TOKEN",
			Message:   "Invalid authentication token",
			RequestID: requestID,
		})

	case errors.Is(err, provider.ErrUserNotFound):
		c.JSON(http.StatusNotFound, APIError{
			Code:      "USER_NOT_FOUND",
			Message:   "User not found",
			RequestID: requestID,
		})

	default:
		// Handle any other errors as internal server errors
		c.JSON(http.StatusInternalServerError, APIError{
			Code:      "INTERNAL_SERVER_ERROR",
			Message:   "An unexpected error occurred",
			RequestID: requestID,
		})
	}
}

// ValidationError represents validation errors
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// HandleValidationError handles validation errors
func HandleValidationError(c *gin.Context, errs []ValidationError) {
	c.JSON(http.StatusBadRequest, gin.H{
		"code":       "VALIDATION_ERROR",
		"message":    "Invalid request parameters",
		"errors":     errs,
		"request_id": GetRequestID(c),
	})
}

// NewValidationError creates a new validation error
func NewValidationError(field, message string) ValidationError {
	return ValidationError{
		Field:   field,
		Message: message,
	}
}
