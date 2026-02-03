package config

import (
	"os"

	"github.com/godotenv/godotenv"
)

type Config struct {
	APIKey  string
	Port    string
	BaseURL string
}

func Load() (*Config, error) {
	godotenv.Load()

	return &Config{
		APIKey:  os.Getenv("API_KEY"),
		Port:    os.Getenv("PORT"),
		BaseURL: os.Getenv("BASE_URL"),
	}, nil
}
