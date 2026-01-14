package main

import "time"

// HTTP handler timeout configurations
const (
	// RequestTimeout is the maximum duration for HTTP handler operations
	RequestTimeout = 5 * time.Second
)

// Token configuration
const (
	// TokenSize is the size in bytes of generated authentication tokens
	TokenSize = 32
)

// Server configuration defaults
const (
	// DefaultPort is the default HTTP server port
	DefaultPort = "8080"

	// DefaultDSN is the default database connection string for development
	DefaultDSN = "postgres://postgres:postgres@localhost:5432/conferenzetech?sslmode=disable"
)

// User roles for conference registration
const (
	RoleAttendee = "attendee"
	RoleSpeaker = "speaker"
	RoleVolunteer = "volunteer"
)

// ValidRoles is a map of all valid conference roles for quick validation
var ValidRoles = map[string]bool{
	RoleAttendee:  true,
	RoleSpeaker:   true,
	RoleVolunteer: true,
}

// IsValidRole checks if the given role is valid
func IsValidRole(role string) bool {
	return ValidRoles[role]
}
