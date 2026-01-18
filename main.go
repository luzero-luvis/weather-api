package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
		log.Println("server is running")
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Get("/weather", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("this is weather api"))
		})
	})

	if err := http.ListenAndServe(":8001", r); err != nil {
		fmt.Println("error starting", err)
	}
}
