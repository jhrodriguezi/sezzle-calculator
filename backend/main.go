package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"sezzle-calculator/backend/handlers"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", handlers.Health)
	mux.HandleFunc("POST /api/calculate", handlers.Calculate)

	var h http.Handler = mux
	h = handlers.WithCORS(h)
	h = handlers.WithLogging(h)

	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      h,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("starting server on :%s", port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server failed: %v", err)
	}
}
