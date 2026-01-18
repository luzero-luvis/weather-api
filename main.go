package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "hello")
		log.Printf("%s", "serving at port 8000")
	})

	if err := http.ListenAndServe(":8000", nil); err != nil {
		fmt.Println("error starting", err)
	}
}
