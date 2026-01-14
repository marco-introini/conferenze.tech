package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/jackc/pgx/v5"
)

// generateToken creates a new random token
func generateToken() (string, error) {
	token := make([]byte, 32)
	_, err := rand.Read(token)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(token), nil
}

// hashToken returns the SHA-256 hex digest of the given token.
// We persist this hash in the database instead of the raw token.
func hashToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// authMiddleware validates the authentication token and adds user ID to context
func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header required", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		tokenHash := hashToken(parts[1])

		token, err := s.db.GetTokenByHash(r.Context(), tokenHash)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}
			log.Printf("Error getting token: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if token.Revoked {
			http.Error(w, "Token revoked", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), UserIDKey, token.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
