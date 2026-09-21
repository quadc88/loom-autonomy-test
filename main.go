package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"loom-bootstrap-test-5/handlers"
	"loom-bootstrap-test-5/middleware"
	"loom-bootstrap-test-5/store"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	s := store.NewInMemoryStore()
	h := handlers.NewTaskHandler(s)

	mux := http.NewServeMux()
	mux.HandleFunc("/tasks", h.HandleTasks)
	mux.HandleFunc("/tasks/", h.HandleTasks)
	mux.HandleFunc("/health", h.HealthCheck)

	// Configure CORS from environment variables
	corsOpts := middleware.DefaultCORSOptions()
	origins := os.Getenv("ALLOWED_ORIGINS")
	if origins != "" {
		corsOpts.AllowedOrigins = strings.Split(origins, ",")
	}

	// Apply middleware chain: Recovery -> Logging -> CORS -> Handler
	handler := middleware.Chain(mux,
		middleware.Recovery,
		middleware.Logging,
		middleware.CORS(corsOpts),
	)

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
