package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
)

// Test validazione JSON per CreateConferenceRequest
func TestConferenceJSONValidation(t *testing.T) {
	tests := []struct {
		name        string
		jsonBody    string
		shouldError bool
	}{
		{
			name:        "Valid JSON",
			jsonBody:    `{"title":"Test","location":"Milano","date":"2026-01-01T00:00:00Z"}`,
			shouldError: false,
		},
		{
			name:        "Invalid JSON",
			jsonBody:    `{"title":"Test", invalid}`,
			shouldError: true,
		},
		{
			name:        "Empty body",
			jsonBody:    ``,
			shouldError: true,
		},
		{
			name:        "Missing required fields",
			jsonBody:    `{"title":"Test"}`,
			shouldError: false, // JSON valid, ma validazione business fallirebbe
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var conf CreateConferenceRequest
			err := json.Unmarshal([]byte(tt.jsonBody), &conf)

			if tt.shouldError && err == nil {
				t.Errorf("Expected error but got none")
			}
			if !tt.shouldError && err != nil {
				t.Errorf("Expected no error but got: %v", err)
			}
		})
	}
}

// Test helper per creare request con auth context
func TestAuthContextHelper(t *testing.T) {
	userID := uuid.New()

	req := httptest.NewRequest("GET", "/api/me", nil)
	ctx := context.WithValue(req.Context(), UserIDKey, userID)
	req = req.WithContext(ctx)

	// Verifica che il context contenga l'user ID
	extractedID, ok := req.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		t.Error("Expected user ID in context but not found")
	}

	if extractedID != userID {
		t.Errorf("Expected user ID %s, got %s", userID, extractedID)
	}
}

// Test conversione response
func TestConferenceResponseConversion(t *testing.T) {
	website := "https://example.com"
	lat := 45.4642
	lon := 9.1900

	response := ConferenceResponse{
		ID:        uuid.New().String(),
		Title:     "Test Conference",
		Date:      "2026-09-15T00:00:00Z",
		Location:  "Milano",
		Website:   &website,
		Latitude:  &lat,
		Longitude: &lon,
	}

	// Verifica che i campi siano popolati correttamente
	if response.Title != "Test Conference" {
		t.Errorf("Expected title 'Test Conference', got '%s'", response.Title)
	}

	if response.Website == nil {
		t.Error("Expected website to be set")
	} else if *response.Website != website {
		t.Errorf("Expected website '%s', got '%s'", website, *response.Website)
	}

	if response.Latitude == nil {
		t.Error("Expected latitude to be set")
	}

	if response.Longitude == nil {
		t.Error("Expected longitude to be set")
	}
}

// Test JSON encoding/decoding
func TestConferenceJSONRoundtrip(t *testing.T) {
	website := "https://gophercon.it"
	lat := 45.4642
	lon := 9.1900

	original := CreateConferenceRequest{
		Title:     "GopherCon 2026",
		Location:  "Milano",
		Date:      "2026-09-15T10:00:00Z",
		Website:   &website,
		Latitude:  &lat,
		Longitude: &lon,
	}

	// Encode to JSON
	jsonData, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Failed to marshal: %v", err)
	}

	// Decode back
	var decoded CreateConferenceRequest
	err = json.Unmarshal(jsonData, &decoded)
	if err != nil {
		t.Fatalf("Failed to unmarshal: %v", err)
	}

	// Verify
	if decoded.Title != original.Title {
		t.Errorf("Title mismatch: expected '%s', got '%s'", original.Title, decoded.Title)
	}

	if decoded.Location != original.Location {
		t.Errorf("Location mismatch: expected '%s', got '%s'", original.Location, decoded.Location)
	}

	if decoded.Date != original.Date {
		t.Errorf("Date mismatch: expected '%s', got '%s'", original.Date, decoded.Date)
	}

	if decoded.Website == nil || *decoded.Website != *original.Website {
		t.Error("Website mismatch after roundtrip")
	}
}

// Test HTTP headers
func TestConferenceResponseHeaders(t *testing.T) {
	rr := httptest.NewRecorder()

	// Simula risposta JSON
	rr.Header().Set("Content-Type", "application/json")

	response := ConferenceResponse{
		ID:       uuid.New().String(),
		Title:    "Test",
		Location: "Milano",
		Date:     "2026-01-01T00:00:00Z",
	}

	json.NewEncoder(rr).Encode(response)

	// Verifica headers
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", ct)
	}

	// Verifica che la response sia valida JSON
	var decoded ConferenceResponse
	if err := json.NewDecoder(rr.Body).Decode(&decoded); err != nil {
		t.Errorf("Response is not valid JSON: %v", err)
	}

	// Verifica i campi
	if decoded.Title != "Test" {
		t.Errorf("Expected title 'Test', got '%s'", decoded.Title)
	}
}

// Test ErrorResponse
func TestErrorResponse(t *testing.T) {
	rr := httptest.NewRecorder()

	errResp := ErrorResponse{
		Error: "Test error message",
	}

	w := rr
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(errResp)

	// Verifica status code
	if rr.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", rr.Code)
	}

	// Decode e verifica
	var decoded ErrorResponse
	if err := json.NewDecoder(rr.Body).Decode(&decoded); err != nil {
		t.Fatalf("Failed to decode error response: %v", err)
	}

	if decoded.Error != "Test error message" {
		t.Errorf("Expected error 'Test error message', got '%s'", decoded.Error)
	}
}

// Test RegisterToConferenceRequest
func TestRegisterToConferenceRequest(t *testing.T) {
	notes := "Looking forward to the event"
	req := RegisterToConferenceRequest{
		ConferenceID: uuid.New().String(),
		Role:         "attendee",
		Notes:        &notes,
		NeedsRide:    true,
		HasCar:       false,
	}

	// Verifica che i campi siano settati
	if req.Role != "attendee" {
		t.Errorf("Expected role 'attendee', got '%s'", req.Role)
	}

	if req.NeedsRide != true {
		t.Error("Expected NeedsRide to be true")
	}

	if req.Notes == nil {
		t.Error("Expected notes to be set")
	}
}

// Helper functions per test
func newAuthRequest(method, url string, body []byte, userID uuid.UUID) *http.Request {
	req := httptest.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	ctx := context.WithValue(req.Context(), UserIDKey, userID)
	return req.WithContext(ctx)
}

func verifyJSONResponse(t *testing.T, rr *httptest.ResponseRecorder, expectedStatus int) {
	if rr.Code != expectedStatus {
		t.Errorf("Expected status %d, got %d. Body: %s",
			expectedStatus, rr.Code, rr.Body.String())
	}

	contentType := rr.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}
}
