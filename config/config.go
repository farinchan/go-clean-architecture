package config

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config holds all configuration for the application
type Config struct {
	App      AppConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	SMTP     SMTPConfig
}

// AppConfig holds application specific configuration
type AppConfig struct {
	Name  string
	Env   string
	Port  string
	Debug bool
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	Timezone string
}

// RedisConfig holds redis configuration
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret      string
	ExpireHours time.Duration
}

// SMTPConfig holds SMTP configuration
type SMTPConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	FromName string
}

// LoadConfig reads configuration from file or environment variables.
func LoadConfig(path string) (*Config, error) {
	viper.SetConfigFile(path)
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		logrus.Warnf("Config file not found, using environment variables: %v", err)
	}

	config := &Config{
		App: AppConfig{
			Name:  viper.GetString("APP_NAME"),
			Env:   viper.GetString("APP_ENV"),
			Port:  viper.GetString("APP_PORT"),
			Debug: viper.GetBool("APP_DEBUG"),
		},
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			DBName:   viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSLMODE"),
			Timezone: viper.GetString("DB_TIMEZONE"),
		},
		Redis: RedisConfig{
			Host:     viper.GetString("REDIS_HOST"),
			Port:     viper.GetString("REDIS_PORT"),
			Password: viper.GetString("REDIS_PASSWORD"),
			DB:       viper.GetInt("REDIS_DB"),
		},
		JWT: JWTConfig{
			Secret:      viper.GetString("JWT_SECRET"),
			ExpireHours: time.Duration(viper.GetInt("JWT_EXPIRE_HOURS")) * time.Hour,
		},
		SMTP: SMTPConfig{
			Host:     viper.GetString("SMTP_HOST"),
			Port:     viper.GetInt("SMTP_PORT"),
			Username: viper.GetString("SMTP_USERNAME"),
			Password: viper.GetString("SMTP_PASSWORD"),
			From:     viper.GetString("SMTP_FROM"),
			FromName: viper.GetString("SMTP_FROM_NAME"),
		},
	}

	return config, nil
}

// GetDSN returns the database connection string
func (d *DatabaseConfig) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		d.Host, d.User, d.Password, d.DBName, d.Port, d.SSLMode, d.Timezone,
	)
}
