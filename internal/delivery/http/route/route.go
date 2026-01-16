package route

import (
	"kopeta-backend/internal/delivery/http/handler"
	"kopeta-backend/internal/delivery/http/middleware"

	"github.com/gin-gonic/gin"
)

type RouteConfig struct {
	Router      *gin.Engine
	UserHandler *handler.UserHandler
}

func SetupRoutes(config *RouteConfig) {
	// Global middleware
	config.Router.Use(middleware.LoggerMiddleware())
	config.Router.Use(middleware.CORSMiddleware())
	config.Router.Use(gin.Recovery())

	// Health check
	config.Router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"message": "Server is running",
		})
	})

	// API v1 routes
	v1 := config.Router.Group("/api/v1")
	{
		// User routes
		users := v1.Group("/users")
		{
			users.POST("", config.UserHandler.Create)
			users.GET("", config.UserHandler.GetAll)
			users.GET("/:id", config.UserHandler.GetByID)
			users.PUT("/:id", config.UserHandler.Update)
			users.DELETE("/:id", config.UserHandler.Delete)
		}
	}
}
