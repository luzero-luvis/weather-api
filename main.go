package main

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"weather-api/internal/cache"
	"weather-api/internal/client"
	"weather-api/internal/config"
	"weather-api/internal/logger"
	"weather-api/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/godotenv/godotenv"
)

func main() {
	godotenv.Load()

	conf, err := config.Load()
	if err != nil {
		slog.Error("faied to load env", "error", err)
		return
	}

	logger.Setup(conf.Env)

	slog.Info("starting weather api sererver",
		"port", conf.Port,
		"env", conf.Env,
	)

	redisClinet, err := cache.NewRedisClient(conf.RedisAddr, conf.RedisPass, conf.RedisDB)
	if err != nil {
		slog.Error("error connectig with redis", "error", err)
		return
	}

	defer redisClinet.Close()

	r := chi.NewRouter()

	r.Use(middleware.LoggingMiddleware)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))

		slog.Debug("health check called")
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/weather", func(w http.ResponseWriter, r *http.Request) {
			// getting city name from url
			city := r.URL.Query().Get("city")

			// if city name is null just say i need ciry name to fetch the data
			if city == "" {

				slog.Warn("weather request missing city parameter",
					"path", r.URL.Path,
				)

				http.Error(w, "city parameteta needs", http.StatusBadRequest)
				return
			}

			cacheKey := "weather:" + strings.ToLower(city)

			var weather client.WeatherResponse

			err := redisClinet.Get(r.Context(), cacheKey, &weather)

			if err == nil {

				slog.Info("fetched the data from cache", "city", city)
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Hit", "cache")

				PrettyJSON, _ := json.MarshalIndent(weather, "", "")
				w.Write(PrettyJSON)
				return
			}

			slog.Debug("mising cache calling the api", "city", city)
			// calling the func that we alreay written in config
			weatherData, err := client.GetWeather(conf.BaseURL, conf.APIKey, city)
			if err != nil {
				slog.Error("failed to get weather data",
					"city", city,
					"error", err,
				)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if err := redisClinet.Set(r.Context(), cacheKey, &weatherData); err != nil {
				slog.Warn("failed to set the cache to redis", "city", city)
			}

			// tell browser that its  json data
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("X-Cache", "MISS")
			Prettyjson, err := json.MarshalIndent(weatherData, " ", "")
			if err != nil {

				slog.Error("failed to encode weather response",
					"city", city,
					"error", err,
				)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			slog.Info("weather request completed successfully",
				"city", city,
			)

			w.Write(Prettyjson)
		})
	})

	slog.Info("server listening", "port", conf.Port)
	if err := http.ListenAndServe(":"+conf.Port, r); err != nil {
		slog.Error("server failed to start", "error", err)
	}
}
