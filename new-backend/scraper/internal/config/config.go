package config

import (
	"os"
	"strconv"
)

type Config struct {
	MongoURI    string
	DBName      string
	OpenAIKey   string
	ProxyURL    string
	LogLevel    string
	RunOnce     bool // if true, run one cycle and exit (for K8s CronJob mode)
	ImportLawId int  // if > 0, import only this single law and exit
}

func Load() *Config {
	importLawId := 0
	if s := os.Getenv("IMPORT_LAW_ID"); s != "" {
		if id, err := strconv.Atoi(s); err == nil {
			importLawId = id
		}
	}

	return &Config{
		MongoURI:    getEnv("MONGODB_URI", "mongodb://localhost:27017"),
		DBName:      getEnv("DB_NAME", "parliamentdb"),
		OpenAIKey:   os.Getenv("OPENAI_API_KEY"),
		ProxyURL:    os.Getenv("CDEP_PROXY_URL"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		RunOnce:     os.Getenv("RUN_ONCE") == "true",
		ImportLawId: importLawId,
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
