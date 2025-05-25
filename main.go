package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	db, err := ConnectToPostgres()

	if err != nil {
		log.Fatalf("Failed to create database connection: %v", err)
	}

	defer db.pool.Close()

	mux.Handle("GET /health", http.HandlerFunc(HealthHandler))

	muxWithMiddleware := AttachGlobalMiddleware(mux, AttachContentTypeJSON)

	http.ListenAndServe(":8080", muxWithMiddleware)
}
