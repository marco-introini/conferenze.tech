package main

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

// responseWriter wraps http.ResponseWriter to capture status code and size
type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

// newResponseWriter creates a new responseWriter
func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		status:         http.StatusOK, // default status
	}
}

// WriteHeader captures the status code
func (rw *responseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

// Write captures the response size
func (rw *responseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

// loggingMiddleware logs HTTP requests with method, path, status, duration, size, and user
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap the response writer to capture status and size
		wrapped := newResponseWriter(w)

		// Get user ID from context if available
		userID := "anonymous"
		if uid, ok := r.Context().Value(UserIDKey).(uuid.UUID); ok {
			userID = uid.String()
		}

		// Process the request
		next.ServeHTTP(wrapped, r)

		// Calculate duration
		duration := time.Since(start)

		// Log the request with all details
		log.Printf(
			"[%s] %s %s | Status: %d | Size: %d bytes | Duration: %v | User: %s",
			r.RemoteAddr,
			r.Method,
			r.URL.Path,
			wrapped.status,
			wrapped.size,
			duration,
			userID,
		)
	})
}
