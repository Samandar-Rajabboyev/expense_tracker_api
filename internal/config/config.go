package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	DatabaseURL string
	Port        string
	JWTSecret   string
	RedisURL    string
	LogLevel    string
}

func Load() (*Config, error) {
	godotenv.Load()

	return &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Port:        os.Getenv("PORT"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		RedisURL:    os.Getenv("REDIS_URL"),
		LogLevel:    os.Getenv("LOG_LEVEL"),
	}, nil
}
