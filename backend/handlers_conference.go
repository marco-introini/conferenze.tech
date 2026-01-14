package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/marco-introini/conferenze.tech/backend/db"
)

// ListConferences retrieves all conferences
func (s *Server) ListConferences(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), RequestTimeout)
	defer cancel()

	conferences, err := s.db.ListConferences(ctx)
	if err != nil {
		log.Printf("Error listing conferences: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := make([]ConferenceResponse, len(conferences))
	for i, c := range conferences {
		dateStr := c.Date.Format(time.RFC3339)
		response[i] = ConferenceResponse{
			ID:        c.ID.String(),
			Title:     c.Title,
			Date:      dateStr,
			Location:  c.Location,
			Website:   stringPtr(c.Website),
			Latitude:  float64Ptr(c.Latitude),
			Longitude: float64Ptr(c.Longitude),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode conferences response: %v", err)
	}
}

// GetConference retrieves a specific conference with its attendees
func (s *Server) GetConference(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), RequestTimeout)
	defer cancel()

	idStr := r.PathValue("conference_id")
	if idStr == "" {
		http.Error(w, "Conference ID required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid conference ID", http.StatusBadRequest)
		return
	}

	conference, err := s.db.GetConferenceByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Conference not found", http.StatusNotFound)
			return
		}
		log.Printf("Error getting conference: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	registrations, err := s.db.GetRegistrationsByConference(ctx, id)
	if err != nil {
		log.Printf("Error getting registrations: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := ConferenceWithAttendees{
		ID:        conference.ID.String(),
		Title:     conference.Title,
		Date:      conference.Date.Format(time.RFC3339),
		Location:  conference.Location,
		Website:   stringPtr(conference.Website),
		Latitude:  float64Ptr(conference.Latitude),
		Longitude: float64Ptr(conference.Longitude),
		Attendees: make([]Attendee, 0),
	}

	for _, reg := range registrations {
		attendee := Attendee{
			User: UserResponse{
				ID:        reg.UserID.String(),
				Email:     reg.Email,
				Name:      reg.Name,
				Nickname:  stringPtr(reg.Nickname),
				City:      stringPtr(reg.City),
				AvatarURL: stringPtr(reg.AvatarUrl),
			},
			NeedsRide: boolPtr(reg.NeedsRide),
			HasCar:    boolPtr(reg.HasCar),
		}
		response.Attendees = append(response.Attendees, attendee)
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode conference response: %v", err)
	}
}

// CreateConference creates a new conference
func (s *Server) CreateConference(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), RequestTimeout)
	defer cancel()

	userID, ok := r.Context().Value(UserIDKey).(uuid.UUID)
	if !ok {
		http.Error(w, "User not found in context", http.StatusNotFound)
		return
	}

	var req CreateConferenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Title == "" || req.Location == "" || req.Date == "" {
		http.Error(w, "Title, location and date are required", http.StatusBadRequest)
		return
	}

	date, err := time.Parse(time.RFC3339, req.Date)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		return
	}

	conference, err := s.db.CreateConference(ctx, db.CreateConferenceParams{
		Title:     req.Title,
		Date:      date,
		Location:  req.Location,
		Website:   nullString(req.Website),
		Latitude:  nullFloat64(req.Latitude),
		Longitude: nullFloat64(req.Longitude),
		CreatedBy: userID,
	})
	if err != nil {
		log.Printf("Error creating conference: %v", err)
		http.Error(w, "Failed to create conference", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(ConferenceResponse{
		ID:        conference.ID.String(),
		Title:     conference.Title,
		Date:      conference.Date.Format(time.RFC3339),
		Location:  conference.Location,
		Website:   stringPtr(conference.Website),
		Latitude:  float64Ptr(conference.Latitude),
		Longitude: float64Ptr(conference.Longitude),
		CreatedBy: userID.String(),
	}); err != nil {
		log.Printf("Failed to encode create conference response: %v", err)
	}
}
