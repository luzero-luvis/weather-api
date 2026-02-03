package main

import (
	"encoding/json"
	"fmt"
	"log"
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
		log.Fatal(err)
	}

	log.Printf("Starting server on port %s", conf.Port)

	r := chi.NewRouter()
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
		log.Println("server is running")
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/weather", func(w http.ResponseWriter, r *http.Request) {
			// getting city name from url
			city := r.URL.Query().Get("city")

			// if city name is null just say i need ciry name to fetch the data
			if city == "" {
				http.Error(w, "city parameteta needs", http.StatusBadRequest)
			}

			// calling the func that we alreay written in config
			weather, err := client.GetWeather(conf.BaseURL, conf.APIKey, city)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			// tell browser that its  json data
			w.Header().Set("Content-Type", "application/json")

			Prettyjson, err := json.MarshalIndent(weather, " ", "")
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			w.Write(Prettyjson)
		})
	})

	if err := http.ListenAndServe(":"+conf.Port, r); err != nil {
		fmt.Println("error starting", err)
	}
}
