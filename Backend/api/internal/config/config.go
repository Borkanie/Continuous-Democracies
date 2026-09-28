// Package config loads the brand_new_day API configuration from environment
// variables. Pattern reused from new-backend/api/internal/config/config.go.
package config

import (
	"log/slog"
	"os"
	"strings"
)

// Configuration holds every tunable the API server and the seeder need.
type Configuration struct {
	MongoURI     string
	DatabaseName string
	Port         string
	LogLevel     string
}

// Load reads the configuration from the environment, applying defaults for any
// variable that is unset or empty.
func Load() *Configuration {
	return &Configuration{
		MongoURI:     readEnvironmentVariable("MONGODB_URI", "mongodb://localhost:27017"),
		DatabaseName: readEnvironmentVariable("DB_NAME", "brandnewdaydb"),
		Port:         readEnvironmentVariable("PORT", "8090"),
		LogLevel:     readEnvironmentVariable("LOG_LEVEL", "info"),
	}
}

// SlogLevel translates the textual LOG_LEVEL into a slog.Level.
func (configuration *Configuration) SlogLevel() slog.Level {
	switch strings.ToLower(configuration.LogLevel) {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func readEnvironmentVariable(environmentVariableName string, fallbackValue string) string {
	if rawValue := os.Getenv(environmentVariableName); rawValue != "" {
		return rawValue
	}
	return fallbackValue
}
