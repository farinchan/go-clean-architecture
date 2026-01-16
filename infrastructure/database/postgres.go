package database

import (
	"go-clean-architecture/config"
	"go-clean-architecture/pkg/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func NewPostgresDB(cfg *config.DatabaseConfig) (*gorm.DB, error) {
	dsn := cfg.GetDSN()

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormlogger.Default.LogMode(gormlogger.Info),
	})
	if err != nil {
		logger.Errorf("Failed to connect to database: %v", err)
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		logger.Errorf("Failed to get database instance: %v", err)
		return nil, err
	}

	// Set connection pool settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	logger.Info("Successfully connected to PostgreSQL database")
	return db, nil
}
