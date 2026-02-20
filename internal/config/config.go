package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
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

	required := map[string]string{
		"APIKey":    conf.APIKey,
		"BaseURL":   conf.BaseURL,
		"RedisAddr": conf.RedisAddr,
		"RedisPass": conf.RedisPass,
	}

	for key, val := range required {
		if val == "" {
			return nil, fmt.Errorf("required env is missing: %s", key)
		}
	}

	slog.Info("configuration loaded",
		"port", conf.Port,
		"env", conf.Env,
		"RedisAddr", conf.RedisAddr,
	)
	return conf, nil
}
