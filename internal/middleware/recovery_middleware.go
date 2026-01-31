package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/your-username/go-clean-architecture/pkg/logger"
	"github.com/your-username/go-clean-architecture/pkg/response"
)

// RecoveryMiddleware creates a recovery middleware that handles panics
func RecoveryMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Errorf("Panic recovered: %v", err)
				response.Error(c, http.StatusInternalServerError, "Internal server error", nil)
				c.Abort()
			}
		}()
		c.Next()
	}
}
