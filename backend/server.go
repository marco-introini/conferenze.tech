package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/marco-introini/conferenze.tech/backend/db"
)

// Server represents the HTTP server with database access
type Server struct {
	db *db.Queries
}

// NewServer creates a new Server instance
func NewServer(database *db.Queries) *Server {
	return &Server{db: database}
}

// Run starts the HTTP server with all routes configured
func (s *Server) Run(port string) error {
	mux := http.NewServeMux()

	// Public routes (no authentication required)
	mux.HandleFunc("POST /api/register", s.Register)
	mux.HandleFunc("POST /api/login", s.Login)

	// Protected routes (authentication required)
	protected := http.NewServeMux()
	protected.HandleFunc("GET /api/conferences", s.ListConferences)
	protected.HandleFunc("GET /api/conferences/{conference_id}", s.GetConference)
	protected.HandleFunc("POST /api/conferences/create", s.CreateConference)
	protected.HandleFunc("POST /api/conferences/{conference_id}/register", s.RegisterToConference)
	protected.HandleFunc("GET /api/registrations/{user_id}", s.GetUserRegistrations)
	protected.HandleFunc("GET /api/users/{user_id}", s.GetMe)
	protected.HandleFunc("GET /api/me", s.GetMeFromToken)
	protected.HandleFunc("GET /api/tokens", s.GetTokens)
	protected.HandleFunc("POST /api/token/revoke", s.RevokeToken)

	// Mount protected routes with authentication middleware
	mux.Handle("/api/", s.authMiddleware(protected))

	// Health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := io.WriteString(w, "OK"); err != nil {
			log.Printf("Failed to write health check response: %v", err)
		}
	})

	// Apply middleware chain
	handler := loggingMiddleware(corsMiddleware(mux))

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting on %s", addr)
	return http.ListenAndServe(addr, handler)
}
