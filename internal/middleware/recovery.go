package middleware

import (
	"fmt"
	"github.com/zahidhasanpapon/iam-bridge/pkg/logger"
	"net/http"
	"runtime/debug"

	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware returns a middleware that recovers from panics
func RecoveryMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// Get stack trace
				stack := debug.Stack()

				// Log the error with structured fields
				log.Error(
					"Panic recovered",
					"error", err,
					"stack", string(stack),
					"request_id", GetRequestID(c),
					"path", c.Request.URL.Path,
					"method", c.Request.Method,
				)

				// Return error response
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"error":      "Internal Server Error",
					"code":       "INTERNAL_SERVER_ERROR",
					"request_id": GetRequestID(c),
				})
			}
		}()

		c.Next()
	}
}

// PanicHandler handles panic in a controlled way
func PanicHandler(err interface{}, c *gin.Context, log logger.Logger) {
	var errMsg string
	switch v := err.(type) {
	case error:
		errMsg = v.Error()
	case string:
		errMsg = v
	default:
		errMsg = fmt.Sprintf("%v", v)
	}

	// Log error with structured fields
	log.Error(
		"Panic occurred",
		"error", errMsg,
		"stack", string(debug.Stack()),
		"request_id", GetRequestID(c),
	)

	c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
		"error":      "Internal Server Error",
		"code":       "INTERNAL_SERVER_ERROR",
		"request_id": GetRequestID(c),
	})
}
