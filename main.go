package main

import (
	"log"
	"net/http"
	"os"

	"loom-bootstrap-test-5/handlers"
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

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, mux))
}
