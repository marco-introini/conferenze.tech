package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/google/uuid"
)

// GetTokens retrieves all tokens for a specific user
func (s *Server) GetTokens(w http.ResponseWriter, r *http.Request) {
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

	tokens, err := s.db.GetTokensByUser(ctx, userID)
	if err != nil {
		log.Printf("Error getting tokens: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	response := make([]TokenResponse, len(tokens))
	for i, t := range tokens {
		response[i] = TokenResponse{
			ID:         t.ID.String(),
			CreatedAt:  *timePtr(t.CreatedAt),
			LastUsedAt: timePtr(t.LastUsedAt),
			Revoked:    t.Revoked,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode tokens response: %v", err)
	}
}

// RevokeToken revokes a specific token
func (s *Server) RevokeToken(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), RequestTimeout)
	defer cancel()

	// Accept POST or DELETE for revocation
	if r.Method != "POST" && r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, "Token ID required", http.StatusBadRequest)
		return
	}

	id, err := uuid.Parse(idStr)
	if err != nil {
		http.Error(w, "Invalid token ID", http.StatusBadRequest)
		return
	}

	token, err := s.db.RevokeToken(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Token not found", http.StatusNotFound)
			return
		}
		log.Printf("Error revoking token: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	resp := TokenResponse{
		ID:         token.ID.String(),
		CreatedAt:  *timePtr(token.CreatedAt),
		LastUsedAt: timePtr(token.LastUsedAt),
		Revoked:    token.Revoked,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("Failed to encode revoke token response: %v", err)
	}
}
