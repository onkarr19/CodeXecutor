package app

import (
	"CodeXecutor/internal/app/handler"
	"CodeXecutor/internal/middleware"
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Server represents the application server.
type Server struct {
	// Add server-related fields and dependencies here
	httpServer *http.Server
}

// NewServer initializes and returns a new Server instance.
func NewServer() *Server {
	return &Server{}
}

// Start starts the application server.
func (server *Server) Start() {
	// Initialize Gorilla mux  router
	router := mux.NewRouter()

	// Use the JSONMiddleware for all routes
	router.Use(middleware.JSONMiddleware)

	// Define routes
	router.HandleFunc("/submit", handler.HandleCodeSubmission).Methods("POST")
	router.HandleFunc("/result", handler.HandleResult).Methods("GET")

	// Create an HTTP server with the Gorilla Mux router
	server.httpServer = &http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}

	go func() {
		if err := server.httpServer.ListenAndServe(); err != nil {
			log.Fatalf("HTTP server error: %v", err)
		}
	}()

	log.Println("Server started on localhost:8080")
}

// Stop gracefully stops the application server.
func (server *Server) Stop() {
	// Shutdown the HTTP server gracefully
	if err := server.httpServer.Shutdown(context.Background()); err != nil {
		log.Printf("Error during server shutdown: %v", err)
	}
	log.Println("shutting the server")
}
