package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/your-username/go-clean-architecture/internal/handler"
	"github.com/your-username/go-clean-architecture/internal/middleware"
	"github.com/your-username/go-clean-architecture/pkg/utils"
)

// Router holds all route configurations
type Router struct {
	engine        *gin.Engine
	userHandler   *handler.UserHandler
	healthHandler *handler.HealthHandler
	jwtManager    *utils.JWTManager
}

// NewRouter creates a new router instance
func NewRouter(
	userHandler *handler.UserHandler,
	healthHandler *handler.HealthHandler,
	jwtManager *utils.JWTManager,
	debug bool,
) *Router {
	if debug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	engine := gin.New()

	return &Router{
		engine:        engine,
		userHandler:   userHandler,
		healthHandler: healthHandler,
		jwtManager:    jwtManager,
	}
}

// SetupRoutes sets up all routes
func (r *Router) SetupRoutes() *gin.Engine {
	// Global middleware
	r.engine.Use(middleware.RecoveryMiddleware())
	r.engine.Use(middleware.LoggerMiddleware())
	r.engine.Use(middleware.CORSMiddleware())

	// Health check routes (no auth required)
	r.engine.GET("/health", r.healthHandler.Health)
	r.engine.GET("/ready", r.healthHandler.Ready)

	// Swagger documentation
	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 routes
	v1 := r.engine.Group("/api/v1")
	{
		// Auth routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/register", r.userHandler.Register)
			auth.POST("/login", r.userHandler.Login)
		}

		// User routes (protected)
		users := v1.Group("/users")
		users.Use(middleware.AuthMiddleware(r.jwtManager))
		{
			users.GET("/me", r.userHandler.GetCurrentUser)
			users.GET("", r.userHandler.GetUsers)
			users.GET("/:id", r.userHandler.GetUser)
			users.PUT("/:id", r.userHandler.UpdateUser)
			users.DELETE("/:id", r.userHandler.DeleteUser)
		}

		// Admin routes (protected with role check)
		admin := v1.Group("/admin")
		admin.Use(middleware.AuthMiddleware(r.jwtManager))
		admin.Use(middleware.RoleMiddleware("admin"))
		{
			// Add admin-only routes here
		}
	}

	return r.engine
}

// GetEngine returns the gin engine
func (r *Router) GetEngine() *gin.Engine {
	return r.engine
}
