package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
}

type AppConfig struct {
	Name string
	Port string
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("env")

	// Set default values
	viper.SetDefault("APP_NAME", "kopeta-backend")
	viper.SetDefault("APP_PORT", "8080")
	viper.SetDefault("APP_ENV", "development")
	viper.SetDefault("DB_HOST", "localhost")
	viper.SetDefault("DB_PORT", "5432")
	viper.SetDefault("DB_USER", "postgres")
	viper.SetDefault("DB_PASSWORD", "postgres")
	viper.SetDefault("DB_NAME", "kopeta_db")
	viper.SetDefault("DB_SSLMODE", "disable")
	viper.SetDefault("REDIS_HOST", "localhost")
	viper.SetDefault("REDIS_PORT", "6379")
	viper.SetDefault("REDIS_PASSWORD", "")
	viper.SetDefault("REDIS_DB", 0)

	// Read from .env file
	if err := viper.ReadInConfig(); err != nil {
		// If .env file not found, continue with env variables
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			fmt.Printf("Warning: Error reading config file: %v\n", err)
		}
	}

	// Enable reading from environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	config := &Config{
		App: AppConfig{
			Name: viper.GetString("APP_NAME"),
			Port: viper.GetString("APP_PORT"),
			Env:  viper.GetString("APP_ENV"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetString("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
	}

	return config, nil
}

func (d *DatabaseConfig) GetDSN() string {
	return "host=" + d.Host +
		" user=" + d.User +
		" password=" + d.Password +
		" dbname=" + d.DBName +
		" port=" + d.Port +
		" sslmode=" + d.SSLMode
}
