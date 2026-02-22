package config

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	APIKey    string
	Port      string
	BaseURL   string
	Env       string
	RedisAddr string
	RedisPass string
	RedisDB   int
	TLL       int
}

func Load() (*Config, error) {
	redisDB, _ := strconv.Atoi(os.Getenv("DB"))
	redisTls, _ := strconv.Atoi(os.Getenv("TLL"))
	conf := &Config{
		APIKey:    os.Getenv("API_KEY"),
		Port:      os.Getenv("PORT"),
		BaseURL:   os.Getenv("BASE_URL"),
		Env:       os.Getenv("ENV"),
		RedisAddr: os.Getenv("ADDR"),
		RedisPass: os.Getenv("REDIS_PASS"),
		RedisDB:   redisDB,
		TLL:       redisTls,
	}

	// im using interface here i can be any type

	required := map[string]interface{}{
		"APIKey":    conf.APIKey,
		"BaseURL":   conf.BaseURL,
		"RedisAddr": conf.RedisAddr,
		"RedisPass": conf.RedisPass,
		"RedisTls":  conf.TLL,
	}
	var missing []string
	for key, val := range required {
		if val == "" {
			missing = append(missing, key)
		}
	}
	if len(missing) > 0 {
		for _, key := range missing {
			slog.Error("missing env", "key", key)
		}
		return nil, fmt.Errorf("missing env(s): %s", strings.Join(missing, ""))
	}

	slog.Info("configuration loaded",
		"port", conf.Port,
		"env", conf.Env,
		"RedisAddr", conf.RedisAddr,
	)
	return conf, nil
}
