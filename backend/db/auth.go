package db

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// HashPassword genera un hash SHA-256 combinato con un salt casuale.
// Ritorna una stringa nel formato salt:hash
func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	hash := sha256.Sum256(append([]byte(password), salt...))
	return hex.EncodeToString(salt) + ":" + hex.EncodeToString(hash[:]), nil
}

// CheckPasswordHash verifica se una password corrisponde all'hash memorizzato.
func CheckPasswordHash(password, stored string) bool {
	parts := strings.Split(stored, ":")
	if len(parts) != 2 {
		return false
	}

	salt, err := hex.DecodeString(parts[0])
	if err != nil {
		return false
	}

	hash := sha256.Sum256(append([]byte(password), salt...))
	return parts[1] == hex.EncodeToString(hash[:])
}
