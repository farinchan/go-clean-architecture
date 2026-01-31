package logger

import (
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

var Log *logrus.Logger

// InitLogger initializes the logrus logger
func InitLogger(debug bool) {
	Log = logrus.New()

	// Set output to stdout
	Log.SetOutput(os.Stdout)

	// Set log level
	if debug {
		Log.SetLevel(logrus.DebugLevel)
	} else {
		Log.SetLevel(logrus.InfoLevel)
	}

	// Set JSON formatter for production, text for development
	if os.Getenv("APP_ENV") == "production" {
		Log.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat: time.RFC3339,
		})
	} else {
		Log.SetFormatter(&logrus.TextFormatter{
			FullTimestamp:   true,
			TimestampFormat: "2006-01-02 15:04:05",
			ForceColors:     true,
		})
	}
}

// Info logs info level message
func Info(args ...interface{}) {
	Log.Info(args...)
}

// Infof logs info level message with format
func Infof(format string, args ...interface{}) {
	Log.Infof(format, args...)
}

// Debug logs debug level message
func Debug(args ...interface{}) {
	Log.Debug(args...)
}

// Debugf logs debug level message with format
func Debugf(format string, args ...interface{}) {
	Log.Debugf(format, args...)
}

// Warn logs warn level message
func Warn(args ...interface{}) {
	Log.Warn(args...)
}

// Warnf logs warn level message with format
func Warnf(format string, args ...interface{}) {
	Log.Warnf(format, args...)
}

// Error logs error level message
func Error(args ...interface{}) {
	Log.Error(args...)
}

// Errorf logs error level message with format
func Errorf(format string, args ...interface{}) {
	Log.Errorf(format, args...)
}

// Fatal logs fatal level message and exits
func Fatal(args ...interface{}) {
	Log.Fatal(args...)
}

// Fatalf logs fatal level message with format and exits
func Fatalf(format string, args ...interface{}) {
	Log.Fatalf(format, args...)
}

// WithField creates a log entry with a field
func WithField(key string, value interface{}) *logrus.Entry {
	return Log.WithField(key, value)
}

// WithFields creates a log entry with multiple fields
func WithFields(fields logrus.Fields) *logrus.Entry {
	return Log.WithFields(fields)
}
