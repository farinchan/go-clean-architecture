package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func main() {
	// Parse command line flags
	direction := flag.String("direction", "up", "Migration direction: up, down, force, version")
	steps := flag.Int("steps", 0, "Number of migrations to run (0 = all)")
	forceVersion := flag.Int("force", -1, "Force migration to a specific version")
	flag.Parse()

	// Load config
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Config file not found, using environment variables: %v", err)
	}

	// Build DSN
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		viper.GetString("DB_USER"),
		viper.GetString("DB_PASSWORD"),
		viper.GetString("DB_HOST"),
		viper.GetString("DB_PORT"),
		viper.GetString("DB_NAME"),
		viper.GetString("DB_SSLMODE"),
	)

	// Connect to database
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Create driver
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create database driver: %v", err)
	}

	// Create migration instance
	m, err := migrate.NewWithDatabaseInstance(
		"file://database/migrations",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}

	// Run migration based on direction
	switch *direction {
	case "up":
		if *steps > 0 {
			err = m.Steps(*steps)
		} else {
			err = m.Up()
		}
	case "down":
		if *steps > 0 {
			err = m.Steps(-*steps)
		} else {
			err = m.Down()
		}
	case "force":
		if *forceVersion >= 0 {
			err = m.Force(*forceVersion)
		} else {
			log.Fatal("Please provide a version number with -force flag")
		}
	case "version":
		version, dirty, verr := m.Version()
		if verr != nil {
			log.Fatalf("Failed to get version: %v", verr)
		}
		fmt.Printf("Current version: %d, Dirty: %v\n", version, dirty)
		os.Exit(0)
	default:
		log.Fatalf("Unknown direction: %s", *direction)
	}

	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}

	fmt.Println("Migration completed successfully!")
}
