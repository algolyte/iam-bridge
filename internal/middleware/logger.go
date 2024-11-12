package middleware

import (
	"bytes"
	"github.com/zahidhasanpapon/iam-bridge/pkg/logger"
	"io"
	"time"

	"github.com/gin-gonic/gin"
)

// responseWriter captures the status code and response size
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// LoggerMiddleware returns a middleware for logging HTTP requests
func LoggerMiddleware(log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Create a copy of the request body for logging
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Create a custom response writer
		w := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = w

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Log request details
		fields := map[string]interface{}{
			"client_ip":     c.ClientIP(),
			"method":        c.Request.Method,
			"path":          c.Request.URL.Path,
			"status":        c.Writer.Status(),
			"duration":      duration.String(),
			"duration_ms":   duration.Milliseconds(),
			"user_agent":    c.Request.UserAgent(),
			"error":         c.Errors.String(),
			"response_size": c.Writer.Size(),
		}

		// Add request body if present and not too large
		if len(requestBody) > 0 && len(requestBody) < 1024 {
			fields["request_body"] = string(requestBody)
		}

		// Add response body if present and not too large
		if w.body.Len() > 0 && w.body.Len() < 1024 {
			fields["response_body"] = w.body.String()
		}

		// Log based on status code
		status := c.Writer.Status()
		switch {
		case status >= 500:
			log.Error("Server error", fields)
		case status >= 400:
			log.Info("Client error", fields)
		case status >= 300:
			log.Info("Redirection", fields)
		default:
			log.Info("Success", fields)
		}
	}
}
