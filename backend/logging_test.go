package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

func TestLoggingMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		path           string
		userID         *uuid.UUID
		expectedStatus int
		responseBody   string
	}{
		{
			name:           "Anonymous GET request",
			method:         "GET",
			path:           "/api/conferences",
			userID:         nil,
			expectedStatus: http.StatusOK,
			responseBody:   "test response",
		},
		{
			name:           "Authenticated POST request",
			method:         "POST",
			path:           "/api/conferences/create",
			userID:         ptr(uuid.New()),
			expectedStatus: http.StatusCreated,
			responseBody:   `{"id":"123"}`,
		},
		{
			name:           "Error response",
			method:         "GET",
			path:           "/api/invalid",
			userID:         nil,
			expectedStatus: http.StatusNotFound,
			responseBody:   "Not Found",
		},
		{
			name:           "Authenticated user with error",
			method:         "DELETE",
			path:           "/api/conferences/123",
			userID:         ptr(uuid.New()),
			expectedStatus: http.StatusForbidden,
			responseBody:   "Forbidden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test handler that returns the expected response
			handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify user context if expected
				if tt.userID != nil {
					userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
					if !ok {
						t.Error("Expected user ID in context but not found")
					} else if userID != *tt.userID {
						t.Errorf("Expected user ID %s, got %s", tt.userID.String(), userID.String())
					}
				}

				w.WriteHeader(tt.expectedStatus)
				w.Write([]byte(tt.responseBody))
			})

			// Wrap with logging middleware
			loggedHandler := loggingMiddleware(handler)

			// Create a test request
			req := httptest.NewRequest(tt.method, tt.path, nil)

			// Add user to context if specified
			if tt.userID != nil {
				ctx := context.WithValue(req.Context(), UserIDKey, *tt.userID)
				req = req.WithContext(ctx)
			}

			// Create a response recorder
			rr := httptest.NewRecorder()

			// Execute the request
			loggedHandler.ServeHTTP(rr, req)

			// Verify status code
			if rr.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, rr.Code)
			}

			// Verify response body
			if rr.Body.String() != tt.responseBody {
				t.Errorf("Expected body %q, got %q", tt.responseBody, rr.Body.String())
			}

			// Verify response size matches
			expectedSize := len(tt.responseBody)
			if rr.Body.Len() != expectedSize {
				t.Errorf("Expected size %d, got %d", expectedSize, rr.Body.Len())
			}
		})
	}
}

func TestResponseWriter(t *testing.T) {
	t.Run("Captures status code", func(t *testing.T) {
		rr := httptest.NewRecorder()
		rw := newResponseWriter(rr)

		rw.WriteHeader(http.StatusCreated)

		if rw.status != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, rw.status)
		}
	})

	t.Run("Default status is 200", func(t *testing.T) {
		rr := httptest.NewRecorder()
		rw := newResponseWriter(rr)

		// Don't call WriteHeader, should default to 200
		rw.Write([]byte("test"))

		if rw.status != http.StatusOK {
			t.Errorf("Expected default status %d, got %d", http.StatusOK, rw.status)
		}
	})

	t.Run("Captures response size", func(t *testing.T) {
		rr := httptest.NewRecorder()
		rw := newResponseWriter(rr)

		testData := []byte("test response body")
		n, err := rw.Write(testData)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}

		if n != len(testData) {
			t.Errorf("Expected to write %d bytes, wrote %d", len(testData), n)
		}

		if rw.size != len(testData) {
			t.Errorf("Expected size %d, got %d", len(testData), rw.size)
		}
	})

	t.Run("Captures multiple writes", func(t *testing.T) {
		rr := httptest.NewRecorder()
		rw := newResponseWriter(rr)

		rw.Write([]byte("first "))
		rw.Write([]byte("second "))
		rw.Write([]byte("third"))

		expectedSize := len("first second third")
		if rw.size != expectedSize {
			t.Errorf("Expected total size %d, got %d", expectedSize, rw.size)
		}
	})
}

// Helper function to create pointer to UUID
func ptr(u uuid.UUID) *uuid.UUID {
	return &u
}
