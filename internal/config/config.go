package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort     string
	OPENAI_API_KEY string
	DEPLOY_SECRET  string
	REDIS_URL      string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables only.")
	}

	return &Config{
		ServerPort:     getEnv("SERVER_PORT", ":8080"),
		OPENAI_API_KEY: getEnv("OPENAI_API_KEY", "sk-********"),
		DEPLOY_SECRET:  getEnv("DEPLOY_SECRET", "secret"),
		REDIS_URL:      getEnv("REDIS_URL", "localhost"),
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
