package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"weather-api/internal/client"
	"weather-api/internal/config"

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

	logger.setUp(conf.Env)

	slog.Info("starting weather api sererver",
		"port", conf.Port,
		"env", conf.Env,
	)

	r := chi.NewRouter()

	r.Use(middleware.loggigMiddleware)

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

			// calling the func that we alreay written in config
			weather, err := client.GetWeather(conf.BaseURL, conf.APIKey, city)
			if err != nil {
				slog.Error("failed to get weather data",
					"city", city,
					"error", err,
				)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// tell browser that its  json data
			w.Header().Set("Content-Type", "application/json")

			Prettyjson, err := json.MarshalIndent(weather, " ", "")
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
