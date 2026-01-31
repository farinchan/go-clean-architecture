package database

import (
	"github.com/your-username/go-clean-architecture/internal/entity"
	"github.com/your-username/go-clean-architecture/pkg/logger"
	"gorm.io/gorm"
)

// Seeder handles database seeding
type Seeder struct {
	db *gorm.DB
}

// NewSeeder creates a new seeder instance
func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{db: db}
}

// Seed runs all seeders
func (s *Seeder) Seed() error {
	logger.Info("Running database seeders...")

	if err := s.seedUsers(); err != nil {
		return err
	}

	logger.Info("Database seeding completed!")
	return nil
}

// seedUsers seeds the users table
func (s *Seeder) seedUsers() error {
	users := []entity.User{
		{
			Name:     "Admin User",
			Email:    "admin@example.com",
			Password: "$2a$10$N9qo8uLOickgx2ZMRZoMye.fVKCBd/h.GqwYY.0mvVxQhVGDtJa7C", // password: password123
			Role:     "admin",
			IsActive: true,
		},
		{
			Name:     "Regular User",
			Email:    "user@example.com",
			Password: "$2a$10$N9qo8uLOickgx2ZMRZoMye.fVKCBd/h.GqwYY.0mvVxQhVGDtJa7C", // password: password123
			Role:     "user",
			IsActive: true,
		},
	}

	for _, user := range users {
		var existing entity.User
		if err := s.db.Where("email = ?", user.Email).First(&existing).Error; err == nil {
			logger.Infof("User %s already exists, skipping...", user.Email)
			continue
		}

		if err := s.db.Create(&user).Error; err != nil {
			return err
		}
		logger.Infof("Created user: %s", user.Email)
	}

	return nil
}
