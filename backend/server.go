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
	mux.HandleFunc("GET /api/conferences", s.ListConferences)
	mux.HandleFunc("GET /api/conferences/{conference_id}", s.GetConference)

	// Protected routes (authentication required)
	s.protectedRoute(mux, "POST /api/conferences", s.CreateConference)
	s.protectedRoute(mux, "DELETE /api/conferences/{conference_id}", s.DeleteConference)
	s.protectedRoute(mux, "POST /api/conferences/{conference_id}/register", s.RegisterToConference)
	s.protectedRoute(mux, "GET /api/users/registrations", s.GetUserRegistrations)
	s.protectedRoute(mux, "DELETE /api/users/registrations/{conference_id}", s.UnregisterFromConference)
	s.protectedRoute(mux, "GET /api/users/{user_id}", s.GetMe)
	s.protectedRoute(mux, "GET /api/me", s.GetMeFromToken)
	s.protectedRoute(mux, "GET /api/tokens", s.GetTokens)
	s.protectedRoute(mux, "POST /api/tokens/revoke", s.RevokeToken)

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

// protectedRoute registers a route that requires authentication
func (s *Server) protectedRoute(mux *http.ServeMux, pattern string, handler http.HandlerFunc) {
	mux.Handle(pattern, s.authMiddleware(handler))
}
