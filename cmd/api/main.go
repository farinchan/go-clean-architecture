package main

import (
	"kopeta-backend/config"
	"kopeta-backend/infrastructure/database"
	"kopeta-backend/internal/delivery/http/handler"
	"kopeta-backend/internal/delivery/http/route"
	"kopeta-backend/internal/repository"
	"kopeta-backend/internal/usecase"
	"kopeta-backend/pkg/logger"
	"kopeta-backend/pkg/validator"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize logger
	logger.InitLogger()

	// Load configuration from .env file
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	// Set log level based on environment
	if cfg.App.Env == "development" {
		logger.SetLogLevel("debug")
	}

	// Initialize database
	db, err := database.NewPostgresDB(&cfg.Database)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}

	// Initialize validator
	validate := validator.NewValidator()

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userUseCase, validate)

	// Setup Gin
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Setup routes
	routeConfig := &route.RouteConfig{
		Router:      router,
		UserHandler: userHandler,
	}
	route.SetupRoutes(routeConfig)

	// Start server
	logger.Infof("Starting server on port %s", cfg.App.Port)
	if err := router.Run(":" + cfg.App.Port); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}
