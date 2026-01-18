package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("welcome"))
		log.Println("server is running")
	})
	if err := http.ListenAndServe(":8001", r); err != nil {
		fmt.Println("error starting", err)
	}
}
