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
	// this part load env

	godotenv.Load()

	conf, err := config.Load()
	if err != nil {
		slog.Error("faied to load env", "error", err)
		return
	}

	// this is where you will set env so it can show logs based on env

	logger.Setup(conf.Env)

	slog.Info("starting weather api sererver",
		"port", conf.Port,
		"env", conf.Env,
	)

	// connectig to reids

	redisClinet, err := cache.NewRedisClient(conf.RedisAddr, conf.RedisPass, conf.RedisDB)
	if err != nil {
		slog.Error("error connectig with redis", "error", err)
		return
	}

	// close the connection when work is done

	defer redisClinet.Close()

	r := chi.NewRouter()

	// this where we will use middleware

	r.Use(middleware.LoggingMiddleware)

	// route to get healthz

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))

		slog.Debug("health check called")
	})

	// route to get actual weather data

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

			// no mater what user request in url with city name like LONDon is convert to london so it no duplicate value

			cacheKey := "weather:" + strings.ToLower(city)

			// this is where you try get the cached data before calling actual api

			var weather client.WeatherResponse

			err := redisClinet.Get(r.Context(), cacheKey, &weather)

			if err == nil {

				slog.Info("fetched the data from cache", "city", city)
				w.Header().Set("Content-Type", "application/json")
				w.Header().Set("Hit", "cache")

				// write the cached data in browser

				PrettyJSON, _ := json.MarshalIndent(weather, "", "")
				w.Write(PrettyJSON)
				return
			}

			slog.Debug("mising cache calling the api", "city", city)

			// calling the func that we alreay written in config
			// calling the actual api

			weatherData, err := client.GetWeather(conf.BaseURL, conf.APIKey, city)
			if err != nil {
				slog.Error("failed to get weather data",
					"city", city,
					"error", err,
				)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// caching the data to redis that we alreay called

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
