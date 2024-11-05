package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/zahidhasanpapon/iam-bridge/internal/config"
	"github.com/zahidhasanpapon/iam-bridge/pkg/logger"
	"time"
)

// RequestID adds a unique request ID to each request
func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("RequestID", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// Logger middleware for logging requests
func Logger(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		// Skip logging for health check endpoints
		if path == "/health" || path == "/metrics" {
			return
		}

		if raw != "" {
			path = path + "?" + raw
		}

		log.Info("%s - [%s] %s %s %d %s %s",
			c.GetString("RequestID"),
			c.ClientIP(),
			c.Request.Method,
			path,
			c.Writer.Status(),
			time.Since(start),
			c.Errors.String(),
		)
	}
}

// CORS middleware
func CORS(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", cfg.CorsAllowedOrigins[0])
		c.Header("Access-Control-Allow-Methods", cfg.CorsAllowedMethods[0])
		c.Header("Access-Control-Allow-Headers", cfg.CorsAllowedHeaders[0])
		c.Header("Access-Control-Max-Age", string(rune(cfg.CorsMaxAge)))

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// ErrorHandler middleware for handling errors
func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			c.JSON(c.Writer.Status(), gin.H{
				"error": err.Error(),
			})
		}
	}
}
