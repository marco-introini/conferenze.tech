package db

import (
	"time"

	"github.com/google/uuid"
)

// RegistrationStatus rappresenta lo stato di una registrazione
type RegistrationStatus string

const (
	RegistrationStatusRegistered RegistrationStatus = "registered"
	RegistrationStatusWaitlist   RegistrationStatus = "waitlist"
	RegistrationStatusCancelled  RegistrationStatus = "cancelled"
	RegistrationStatusAttended   RegistrationStatus = "attended"
)

// UserRole rappresenta il ruolo di un utente a una conferenza
type UserRole string

const (
	RoleAttendee  UserRole = "attendee"
	RoleOrganizer UserRole = "organizer"
	RoleSpeaker   UserRole = "speaker"
	RoleVolunteer UserRole = "volunteer"
)

// User rappresenta un utente nel sistema
type User struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	Name      string    `json:"name" db:"name"`
	Nickname  *string   `json:"nickname,omitempty" db:"nickname"`
	City      *string   `json:"city,omitempty" db:"city"`
	AvatarURL *string   `json:"avatarUrl,omitempty" db:"avatar_url"`
	Bio       *string   `json:"bio,omitempty" db:"bio"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// Conference rappresenta una conferenza tech
type Conference struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Date      time.Time `json:"date" db:"date"`
	Location  string    `json:"location" db:"location"`
	Website   *string   `json:"website,omitempty" db:"website"`
	Latitude  *float64  `json:"latitude,omitempty" db:"latitude"`
	Longitude *float64  `json:"longitude,omitempty" db:"longitude"`
	CreatedAt time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt time.Time `json:"updatedAt" db:"updated_at"`
}

// ConferenceRegistration rappresenta la registrazione di un utente a una conferenza
type ConferenceRegistration struct {
	ID           uuid.UUID          `json:"id" db:"id"`
	UserID       uuid.UUID          `json:"userId" db:"user_id"`
	ConferenceID uuid.UUID          `json:"conferenceId" db:"conference_id"`
	Status       RegistrationStatus `json:"status" db:"status"`
	Role         UserRole           `json:"role" db:"role"`
	Notes        *string            `json:"notes,omitempty" db:"notes"`
	NeedsRide    bool               `json:"needsRide" db:"needs_ride"`
	HasCar       bool               `json:"hasCar" db:"has_car"`
	RegisteredAt time.Time          `json:"registeredAt" db:"registered_at"`
	CancelledAt  *time.Time         `json:"cancelledAt,omitempty" db:"cancelled_at"`
}

// ConferenceWithUser include i dettagli dell'utente per le registrazioni
type ConferenceWithUser struct {
	ConferenceRegistration
	Email     string  `json:"email" db:"email"`
	UserName  string  `json:"userName" db:"name"`
	Nickname  *string `json:"nickname,omitempty" db:"nickname"`
	UserCity  *string `json:"userCity,omitempty" db:"city"`
	AvatarURL *string `json:"avatarUrl,omitempty" db:"avatar_url"`
}

// RegistrationWithConference include i dettagli della conferenza per le registrazioni utente
type RegistrationWithConference struct {
	ConferenceRegistration
	Title    string    `json:"title" db:"title"`
	ConfDate time.Time `json:"conferenceDate" db:"date"`
	Location string    `json:"location" db:"location"`
	Website  *string   `json:"website,omitempty" db:"website"`
}

// ConferenceStats contiene le statistiche per una conferenza
type ConferenceStats struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	Title              string    `json:"title" db:"title"`
	TotalRegistrations int64     `json:"totalRegistrations" db:"total_registrations"`
	ConfirmedCount     int64     `json:"confirmedCount" db:"confirmed_count"`
	NeedingRideCount   int64     `json:"needingRideCount" db:"needing_ride_count"`
	OfferingRideCount  int64     `json:"offeringRideCount" db:"offering_ride_count"`
}

// CarpoolUser viene usato per le liste di carpooling
type CarpoolUser struct {
	User
	ConferenceTitle    string  `json:"conferenceTitle" db:"conference_title"`
	ConferenceLocation string  `json:"conferenceLocation" db:"conference_location"`
	Notes              *string `json:"notes,omitempty" db:"notes"`
}
