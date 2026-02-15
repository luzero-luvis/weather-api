package config

import (
	"log/slog"
	"os"

	"github.com/godotenv/godotenv"
)

type Config struct {
	APIKey  string
	Port    string
	BaseURL string
	Env     string
}

func Load() (*Config, error) {
	godotenv.Load()

	conf := &Config{
		APIKey:  os.Getenv("API_KEY"),
		Port:    os.Getenv("PORT"),
		BaseURL: os.Getenv("BASE_URL"),
		Env:     os.Getenv("ENV"),
	}

	slog.Info("configuration loaded",
		"port", conf.Port,
		"env", conf.Env,
	)
	return conf, nil
}
