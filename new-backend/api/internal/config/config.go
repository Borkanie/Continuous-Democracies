package config

import (
	"os"
)

type Config struct {
	MongoURI string
	DBName   string
	Port     string
	LogLevel string
}

func Load() *Config {
	return &Config{
		MongoURI: getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		DBName:   getEnv("DB_NAME", "parliamentdb"),
		Port:     getEnv("PORT", "8080"),
		LogLevel: getEnv("LOG_LEVEL", "warning"),
	}
}

func getEnv(key, fallback string) string {
	if rawValue := os.Getenv(key); rawValue != "" {
		return rawValue
	}
	return fallback
}
