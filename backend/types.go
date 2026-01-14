package main

import "github.com/marco-introini/conferenze.tech/backend/db"

// Context keys
type contextKey string

const UserIDKey contextKey = "userID"

// Request types
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email     string  `json:"email"`
	Password  string  `json:"password"`
	Name      string  `json:"name"`
	Nickname  *string `json:"nickname"`
	City      *string `json:"city"`
	AvatarURL *string `json:"avatarUrl"`
	Bio       *string `json:"bio"`
}

type CreateConferenceRequest struct {
	Title     string   `json:"title"`
	Date      string   `json:"date"`
	Location  string   `json:"location"`
	Website   *string  `json:"website"`
	Latitude  *float64 `json:"latitude"`
	Longitude *float64 `json:"longitude"`
}

type RegisterToConferenceRequest struct {
	ConferenceID string  `json:"conferenceId"`
	Role         string  `json:"role"`
	Notes        *string `json:"notes"`
	NeedsRide    bool    `json:"needsRide"`
	HasCar       bool    `json:"hasCar"`
}

// Response types
type ErrorResponse struct {
	Error string `json:"error"`
}

type LoginResponse struct {
	User  db.User `json:"user"`
	Token string  `json:"token"`
}

type RegisterResponse struct {
	User  db.User `json:"user"`
	Token string  `json:"token"`
}

type ConferenceResponse struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Date      string   `json:"date"`
	Location  string   `json:"location"`
	Website   *string  `json:"website,omitempty"`
	Latitude  *float64 `json:"latitude,omitempty"`
	Longitude *float64 `json:"longitude,omitempty"`
}

type ConferenceWithAttendees struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Date      string     `json:"date"`
	Location  string     `json:"location"`
	Website   *string    `json:"website,omitempty"`
	Latitude  *float64   `json:"latitude,omitempty"`
	Longitude *float64   `json:"longitude,omitempty"`
	Attendees []Attendee `json:"attendees"`
}

type Attendee struct {
	User      UserResponse `json:"user"`
	NeedsRide *bool        `json:"needsRide"`
	HasCar    *bool        `json:"hasCar"`
}

type UserResponse struct {
	ID       string  `json:"id"`
	Nickname *string `json:"nickname,omitempty"`
	City     *string `json:"city,omitempty"`
}

type RegistrationResponse struct {
	ID                 string `json:"id"`
	ConferenceID       string `json:"conferenceId"`
	ConferenceTitle    string `json:"conferenceTitle"`
	ConferenceDate     string `json:"conferenceDate"`
	ConferenceLocation string `json:"conferenceLocation"`
	Status             string `json:"status"`
	Role               string `json:"role"`
	NeedsRide          *bool  `json:"needsRide"`
	HasCar             *bool  `json:"hasCar"`
	RegisteredAt       string `json:"registeredAt"`
}

type TokenResponse struct {
	ID         string  `json:"id"`
	CreatedAt  string  `json:"createdAt"`
	LastUsedAt *string `json:"lastUsedAt,omitempty"`
	Revoked    bool    `json:"revoked"`
}
