package utils

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
	"strings"
	"time"

	"github.com/google/uuid"
)

// GenerateUUID generates a new UUID
func GenerateUUID() string {
	return uuid.New().String()
}

// GenerateRandomString generates a random string of given length
func GenerateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)

	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}

// GenerateOTP generates a random OTP of given length (numeric only)
func GenerateOTP(length int) (string, error) {
	const charset = "0123456789"
	result := make([]byte, length)

	for i := range result {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result[i] = charset[num.Int64()]
	}

	return string(result), nil
}

// GenerateToken generates a random token (URL-safe base64)
func GenerateToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(bytes), nil
}

// SlugFromString creates a URL-friendly slug from a string
func SlugFromString(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "-")
	return s
}

// ParseDateString parses a date string to time.Time
func ParseDateString(dateStr string) (time.Time, error) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02 15:04:05",
		"2006-01-02",
		"02-01-2006",
		"02/01/2006",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, dateStr); err == nil {
			return t, nil
		}
	}

	return time.Time{}, nil
}

// FormatDate formats time.Time to string
func FormatDate(t time.Time, format string) string {
	if format == "" {
		format = "2006-01-02 15:04:05"
	}
	return t.Format(format)
}
