package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/your-username/go-clean-architecture/config"
	_ "github.com/your-username/go-clean-architecture/docs"
	"github.com/your-username/go-clean-architecture/internal/handler"
	"github.com/your-username/go-clean-architecture/internal/repository"
	"github.com/your-username/go-clean-architecture/internal/router"
	"github.com/your-username/go-clean-architecture/internal/usecase"
	"github.com/your-username/go-clean-architecture/pkg/database"
	"github.com/your-username/go-clean-architecture/pkg/logger"
	"github.com/your-username/go-clean-architecture/pkg/utils"
	"github.com/your-username/go-clean-architecture/pkg/validator"
)

// @title Go Clean Architecture API
// @version 1.0
// @description A RESTful API with Go Clean Architecture

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Initialize logger
	logger.InitLogger(true)
	logger.Info("Starting application...")

	// Load configuration
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	// Initialize logger with config
	logger.InitLogger(cfg.App.Debug)

	// Register custom validator
	validator.RegisterGinValidator()

	// Connect to database
	db, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Connect to Redis
	redis, err := database.NewRedisClient(&cfg.Redis)
	if err != nil {
		logger.Warnf("Failed to connect to Redis: %v", err)
		// Continue without Redis, it's optional
	} else {
		defer redis.Close()
	}

	// Initialize JWT Manager
	jwtManager := utils.NewJWTManager(cfg.JWT.Secret, cfg.JWT.ExpireHours)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db.DB)

	// Initialize use cases
	userUseCase := usecase.NewUserUseCase(userRepo, jwtManager)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userUseCase)
	healthHandler := handler.NewHealthHandler()

	// Initialize router
	r := router.NewRouter(userHandler, healthHandler, jwtManager, cfg.App.Debug)
	engine := r.SetupRoutes()

	// Create HTTP server
	server := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      engine,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine
	go func() {
		logger.Infof("Server is running on port %s", cfg.App.Port)
		logger.Infof("Swagger documentation available at http://localhost:%s/swagger/index.html", cfg.App.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited properly")
}
