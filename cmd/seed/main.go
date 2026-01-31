package main

import (
	"github.com/your-username/go-clean-architecture/config"
	dbseeder "github.com/your-username/go-clean-architecture/database/seeder"
	"github.com/your-username/go-clean-architecture/internal/entity"
	"github.com/your-username/go-clean-architecture/pkg/database"
	"github.com/your-username/go-clean-architecture/pkg/logger"
)

func main() {
	// Initialize logger
	logger.InitLogger(true)
	logger.Info("Starting database seeder...")

	// Load configuration
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		logger.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	db, err := database.NewDatabase(&cfg.Database)
	if err != nil {
		logger.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Auto migrate
	if err := db.AutoMigrate(&entity.User{}); err != nil {
		logger.Fatalf("Failed to auto migrate: %v", err)
	}

	// Run seeder
	s := dbseeder.NewSeeder(db.DB)
	if err := s.Seed(); err != nil {
		logger.Fatalf("Failed to seed database: %v", err)
	}

	logger.Info("Database seeding completed successfully!")
}
