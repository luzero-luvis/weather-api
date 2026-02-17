package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/godotenv/godotenv"
)

type Config struct {
	APIKey    string
	Port      string
	BaseURL   string
	Env       string
	RedisAddr string
	RedisPass string
	RedisDB   int
}

func Load() (*Config, error) {
	godotenv.Load()

	redisDB, _ := strconv.Atoi(os.Getenv("DB"))
	conf := &Config{
		APIKey:    os.Getenv("API_KEY"),
		Port:      os.Getenv("PORT"),
		BaseURL:   os.Getenv("BASE_URL"),
		Env:       os.Getenv("ENV"),
		RedisAddr: os.Getenv("ADDR"),
		RedisPass: os.Getenv("REDIS_PASS"),
		RedisDB:   redisDB,
	}

	slog.Info("configuration loaded",
		"port", conf.Port,
		"env", conf.Env,
		"RedisAddr", conf.RedisAddr,
	)
	return conf, nil
}
