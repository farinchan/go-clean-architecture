package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/your-username/go-clean-architecture/pkg/response"
)

// HealthHandler handles health check requests
type HealthHandler struct{}

// NewHealthHandler creates a new health handler
func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

// Health godoc
// @Summary Health check
// @Description Check if the service is running
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	response.Success(c, "Service is running", gin.H{
		"status": "healthy",
	})
}

// Ready godoc
// @Summary Readiness check
// @Description Check if the service is ready to receive traffic
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} response.Response
// @Router /ready [get]
func (h *HealthHandler) Ready(c *gin.Context) {
	response.Success(c, "Service is ready", gin.H{
		"status": "ready",
	})
}
