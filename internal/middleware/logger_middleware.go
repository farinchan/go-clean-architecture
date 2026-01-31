package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/your-username/go-clean-architecture/pkg/logger"
)

// LoggerMiddleware creates a logging middleware
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		startTime := time.Now()

		// Process request
		c.Next()

		// Calculate latency
		latency := time.Since(startTime)

		// Get status code
		statusCode := c.Writer.Status()

		// Get client IP
		clientIP := c.ClientIP()

		// Get request method
		method := c.Request.Method

		// Get request path
		path := c.Request.URL.Path

		// Log format
		entry := logger.WithFields(logrus.Fields{
			"status_code": statusCode,
			"latency":     latency,
			"client_ip":   clientIP,
			"method":      method,
			"path":        path,
		})

		if len(c.Errors) > 0 {
			entry.Error(c.Errors.ByType(gin.ErrorTypePrivate).String())
		} else {
			if statusCode >= 500 {
				entry.Error("Server error")
			} else if statusCode >= 400 {
				entry.Warn("Client error")
			} else {
				entry.Info("Request completed")
			}
		}
	}
}
