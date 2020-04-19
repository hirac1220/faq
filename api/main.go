package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors" // https://github.com/rs/cors
)

func main() {
	// Routing
	r := mux.NewRouter()
	// http://localhost:15000/faq
	r.HandleFunc("/faq", post).Methods("POST")
	// http://localhost:15000/faq
	r.HandleFunc("/faq", get).Methods("GET")

	// Set cors
	c := cors.New(cors.Options{
		// AllowedOrigins: []string{"http://localhost:8001"},
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	})

	log.Fatal(http.ListenAndServe(":8080", c.Handler(r)))
}
