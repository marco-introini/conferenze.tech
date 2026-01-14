package main

import "github.com/marco-introini/conferenze.tech/backend/db"

// contextKey is a custom type for context keys to avoid collisions with other packages
type contextKey string

// UserIDKey is the context key for storing the authenticated user's ID
const UserIDKey contextKey = "userID"

// LoginRequest represents the payload for user authentication.
// Both Email and Password are required fields.
type LoginRequest struct {
	Email    string `json:"email"`    // User's email address
	Password string `json:"password"` // User's password (plain text, will be hashed server-side)
}

// RegisterRequest represents the payload for new user registration.
// Email, Password, and Name are required fields.
// Other fields are optional and can be provided to enrich the user profile.
type RegisterRequest struct {
	Email     string  `json:"email"`     // User's email address (required, must be unique)
	Password  string  `json:"password"`  // User's password (required, min 8 characters recommended)
	Name      string  `json:"name"`      // User's full name (required)
	Nickname  *string `json:"nickname"`  // Optional display nickname for the user
	City      *string `json:"city"`      // Optional city location of the user
	AvatarURL *string `json:"avatarUrl"` // Optional URL to user's avatar image
	Bio       *string `json:"bio"`       // Optional user biography or description
}

// CreateConferenceRequest represents the payload for creating a new conference.
// Title, Date, and Location are required fields.
type CreateConferenceRequest struct {
	Title     string   `json:"title"`     // Conference title (required)
	Date      string   `json:"date"`      // Conference date in RFC3339 format (required)
	Location  string   `json:"location"`  // Conference location, city and country (required)
	Website   *string  `json:"website"`   // Optional conference website URL
	Latitude  *float64 `json:"latitude"`  // Optional GPS latitude coordinate
	Longitude *float64 `json:"longitude"` // Optional GPS longitude coordinate
}

// RegisterToConferenceRequest represents the payload for registering a user to a conference.
// ConferenceID and Role are required fields.
type RegisterToConferenceRequest struct {
	ConferenceID string  `json:"conferenceId"` // UUID of the conference to register to
	Role         string  `json:"role"`         // User's role: "attendee", "speaker", or "volunteer"
	Notes        *string `json:"notes"`        // Optional notes about the registration
	NeedsRide    bool    `json:"needsRide"`    // Whether user needs transportation to the conference
	HasCar       bool    `json:"hasCar"`       // Whether user can provide transportation to others
}

// ErrorResponse represents a standard API error response.
// It provides a human-readable error message to the client.
type ErrorResponse struct {
	Error string `json:"error"` // Human-readable error message
}

// LoginResponse is returned after successful user authentication.
// It contains the user's details (with password removed) and an authentication token.
type LoginResponse struct {
	User  db.User `json:"user"`  // User details (password field is cleared)
	Token string  `json:"token"` // Authentication bearer token for subsequent requests
}

// RegisterResponse is returned after successful user registration.
// It contains the newly created user's details and an authentication token.
type RegisterResponse struct {
	User  db.User `json:"user"`  // Newly created user details (password field is cleared)
	Token string  `json:"token"` // Authentication bearer token for subsequent requests
}

// ConferenceResponse represents a conference in API responses.
// This is the basic conference information returned in lists and single conference views.
type ConferenceResponse struct {
	ID        string   `json:"id"`                  // Conference UUID
	Title     string   `json:"title"`               // Conference title
	Date      string   `json:"date"`                // Conference date in RFC3339 format
	Location  string   `json:"location"`            // Conference location
	Website   *string  `json:"website,omitempty"`   // Optional conference website URL
	Latitude  *float64 `json:"latitude,omitempty"`  // Optional GPS latitude coordinate
	Longitude *float64 `json:"longitude,omitempty"` // Optional GPS longitude coordinate
}

// ConferenceWithAttendees represents a conference with its full list of registered participants.
// This extended response includes all attendee information for detailed conference views.
type ConferenceWithAttendees struct {
	ID        string     `json:"id"`                  // Conference UUID
	Title     string     `json:"title"`               // Conference title
	Date      string     `json:"date"`                // Conference date in RFC3339 format
	Location  string     `json:"location"`            // Conference location
	Website   *string    `json:"website,omitempty"`   // Optional conference website URL
	Latitude  *float64   `json:"latitude,omitempty"`  // Optional GPS latitude coordinate
	Longitude *float64   `json:"longitude,omitempty"` // Optional GPS longitude coordinate
	Attendees []Attendee `json:"attendees"`           // List of all registered attendees
}

// Attendee represents a conference participant with their basic information.
// This includes public user data and transportation preferences.
type Attendee struct {
	User      UserResponse `json:"user"`      // Public user information
	NeedsRide *bool        `json:"needsRide"` // Whether attendee needs transportation
	HasCar    *bool        `json:"hasCar"`    // Whether attendee can provide transportation
}

// UserResponse represents public user information in API responses.
// Only non-sensitive user data is included for privacy.
type UserResponse struct {
	ID       string  `json:"id"`                 // User UUID
	Nickname *string `json:"nickname,omitempty"` // Optional display nickname
	City     *string `json:"city,omitempty"`     // Optional city location
}

// RegistrationResponse represents a conference registration in API responses.
// It includes both registration details and associated conference information.
type RegistrationResponse struct {
	ID                 string `json:"id"`                 // Registration UUID
	ConferenceID       string `json:"conferenceId"`       // Conference UUID
	ConferenceTitle    string `json:"conferenceTitle"`    // Conference title for convenience
	ConferenceDate     string `json:"conferenceDate"`     // Conference date in RFC3339 format
	ConferenceLocation string `json:"conferenceLocation"` // Conference location
	Status             string `json:"status"`             // Registration status (e.g., "confirmed", "pending")
	Role               string `json:"role"`               // User's role at the conference
	NeedsRide          *bool  `json:"needsRide"`          // Whether user needs transportation
	HasCar             *bool  `json:"hasCar"`             // Whether user can provide transportation
	RegisteredAt       string `json:"registeredAt"`       // Registration timestamp in RFC3339 format
}

// TokenResponse represents an authentication token in API responses.
// This is used for listing user's active tokens and managing token lifecycle.
type TokenResponse struct {
	ID         string  `json:"id"`                   // Token UUID
	CreatedAt  string  `json:"createdAt"`            // Token creation timestamp in RFC3339 format
	LastUsedAt *string `json:"lastUsedAt,omitempty"` // Last usage timestamp in RFC3339 format (null if never used)
	Revoked    bool    `json:"revoked"`              // Whether the token has been revoked
}
