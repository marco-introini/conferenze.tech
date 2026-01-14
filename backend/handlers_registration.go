package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/marco-introini/conferenze.tech/backend/db"
)

// RegisterToConference handles user registration to a conference
func (s *Server) RegisterToConference(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), RequestTimeout)
	defer cancel()

	var req RegisterToConferenceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	userIDStr := r.URL.Query().Get("userId")
	if userIDStr == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil || userID == uuid.Nil {
		http.Error(w, "Valid user ID required", http.StatusBadRequest)
		return
	}

	idStr := r.PathValue("conference_id")
	conferenceID, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid conference ID", http.StatusBadRequest)
		return
	}

	_, err = s.db.GetConferenceByID(ctx, conferenceID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			http.Error(w, "Conference not found", http.StatusNotFound)
			return
		}
		log.Printf("Error getting conference: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	role := req.Role
	if !IsValidRole(role) {
		role = RoleAttendee
	}

	registration, err := s.db.RegisterUserToConference(ctx, db.RegisterUserToConferenceParams{
		UserID:       userID,
		ConferenceID: conferenceID,
		Role:         role,
		Notes:        nullString(req.Notes),
		NeedsRide:    nullBool(req.NeedsRide),
		HasCar:       nullBool(req.HasCar),
	})
	if err != nil {
		log.Printf("Error registering user: %v", err)
		http.Error(w, "Failed to register to conference", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(registration); err != nil {
		log.Printf("Failed to encode registration response: %v", err)
	}
}

// GetUserRegistrations retrieves all registrations for a specific user
func (s *Server) GetUserRegistrations(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), RequestTimeout)
	defer cancel()

	userIDStr := r.URL.Query().Get("userId")
	if userIDStr == "" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	registrations, err := s.db.GetRegistrationsByUser(ctx, userID)
	if err != nil {
		log.Printf("Error getting registrations: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := make([]RegistrationResponse, len(registrations))
	for i, reg := range registrations {
		response[i] = RegistrationResponse{
			ID:                 reg.ID.String(),
			ConferenceID:       reg.ConferenceID.String(),
			ConferenceTitle:    reg.Title,
			ConferenceDate:     reg.Date.Format(time.RFC3339),
			ConferenceLocation: reg.Location,
			Status:             string(reg.Status),
			Role:               string(reg.Role),
			NeedsRide:          boolPtr(reg.NeedsRide),
			HasCar:             boolPtr(reg.HasCar),
			RegisteredAt:       *timePtr(reg.RegisteredAt),
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode registrations response: %v", err)
	}
}
